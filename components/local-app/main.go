// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package main

import (
	_ "embed"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	gitpod "github.com/gitpod-io/gitpod/gitpod-protocol"
	appapi "github.com/gitpod-io/gitpod/local-app/api"
	"github.com/gitpod-io/local-app/pkg/auth"
	"github.com/gitpod-io/local-app/pkg/bastion"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	keyring "github.com/zalando/go-keyring"
	"google.golang.org/grpc"
)

var (
	// Version : current version
	Version string = strings.TrimSpace(version)
	//go:embed version.txt
	version string
)

func main() {
	sshConfig := os.Getenv("GITPOD_LCA_SSH_CONFIG")
	if sshConfig == "" {
		if runtime.GOOS == "windows" {
			sshConfig = filepath.Join(os.TempDir(), "gitpod_ssh_config")
		} else {
			sshConfig = filepath.Join("/tmp", "gitpod_ssh_config")
		}
	}

	app := cli.App{
		Name:                 "gitpod-local-companion",
		Usage:                "connect your Gitpod workspaces",
		Action:               DefaultCommand("run"),
		EnableBashCompletion: true,
		Version:              Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "gitpod-host",
				Usage: "URL of the Gitpod installation to connect to",
				EnvVars: []string{
					"GITPOD_HOST",
				},
				Value: "https://gitpod.io",
			},
			&cli.BoolFlag{
				Name:  "mock-keyring",
				Usage: "Don't use system native keyring, but store Gitpod token in memory",
			},
			&cli.BoolFlag{
				Name:  "allow-cors-from-port",
				Usage: "Allow CORS requests from workspace port location",
			},
			&cli.IntFlag{
				Name:  "api-port",
				Usage: "Local App API endpoint's port",
				EnvVars: []string{
					"GITPOD_LCA_API_PORT",
				},
				Value: 63100,
			},
			&cli.BoolFlag{
				Name:  "auto-tunnel",
				Usage: "Enable auto tunneling",
				EnvVars: []string{
					"GITPOD_LCA_AUTO_TUNNEL",
				},
				Value: true,
			},
			&cli.StringFlag{
				Name: "auth-redirect-url",
				EnvVars: []string{
					"GITPOD_LCA_AUTH_REDIRECT_URL",
				},
			},
			&cli.StringFlag{
				Name: "log-level",
				EnvVars: []string{
					"GITPOD_LCA_LOG_LEVEL",
				},
				Value: logrus.InfoLevel.String(),
			},
		},
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(c *cli.Context) error {
					if c.Bool("mock-keyring") {
						keyring.MockInit()
					}
					logLevel, err := logrus.ParseLevel(c.String("log-level"))
					if err != nil {
						logLevel = logrus.InfoLevel
						logrus.WithError(err).WithField("logLevel", c.String("log-level")).Error("failed to parse log-level")
					}
					return run(c.String("gitpod-host"), c.String("ssh_config"), c.Int("api-port"), c.Bool("allow-cors-from-port"), c.Bool("auto-tunnel"), c.String("auth-redirect-url"), logLevel)
				},
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "ssh_config",
						Usage: "produce and update an OpenSSH compatible ssh_config file (defaults to $GITPOD_LCA_SSH_CONFIG)",
						Value: sshConfig,
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func DefaultCommand(name string) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return ctx.App.Command(name).Run(ctx)
	}
}

func run(origin, sshConfig string, apiPort int, allowCORSFromPort bool, autoTunnel bool, authRedirectUrl string, logLevel logrus.Level) error {
	logrus.SetLevel(logLevel)
	logrus.WithField("ssh_config", sshConfig).Info("writing workspace ssh_config file")

	// Trailing slash(es) result in connection issues, so remove them preemptively
	origin = strings.TrimRight(origin, "/")
	originURL, err := url.Parse(origin)
	if err != nil {
		return err
	}
	wsHostRegex := "(\\.[^.]+)\\." + strings.ReplaceAll(originURL.Host, ".", "\\.")
	wsHostRegex = "([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}|[0-9a-z]{2,16}-[0-9a-z]{2,16}-[0-9a-z]{8})" + wsHostRegex
	if allowCORSFromPort {
		wsHostRegex = "([0-9]+)-" + wsHostRegex
	}
	hostRegex, err := regexp.Compile("^" + wsHostRegex)
	if err != nil {
		return err
	}

	var b *bastion.Bastion

	client, err := connectToServer(auth.LoginOpts{GitpodURL: origin, RedirectURL: authRedirectUrl}, func() {
		if b != nil {
			b.FullUpdate()
		}
	})
	if err != nil {
		return err
	}

	cb := bastion.CompositeCallbacks{
		&logCallbacks{},
	}
	s := &bastion.SSHConfigWritingCallback{Path: sshConfig}
	if sshConfig != "" {
		cb = append(cb, s)
	}

	b = bastion.New(client, cb)
	b.EnableAutoTunnel = autoTunnel
	grpcServer := grpc.NewServer()
	appapi.RegisterLocalAppServer(grpcServer, bastion.NewLocalAppService(b, s))
	allowOrigin := func(origin string) bool {
		// Is the origin a subdomain of the installations hostname?
		return hostRegex.Match([]byte(origin))
	}
	go func() {
		err := http.ListenAndServe("localhost:"+strconv.Itoa(apiPort), grpcweb.WrapServer(grpcServer,
			grpcweb.WithCorsForRegisteredEndpointsOnly(false),
			grpcweb.WithOriginFunc(allowOrigin),
			grpcweb.WithWebsockets(true),
			grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
				origin, err := grpcweb.WebsocketRequestOrigin(req)
				if err != nil {
					return false
				}
				return allowOrigin(origin)
			}),
			grpcweb.WithWebsocketPingInterval(15*time.Second),
		))
		if err != nil {
			logrus.WithError(err).Error("API endpoint failed to start")
			os.Exit(1)
		}
	}()
	defer grpcServer.Stop()
	return b.Run()
}

func connectToServer(loginOpts auth.LoginOpts, reconnectionHandler func()) (*gitpod.APIoverJSONRPC, error) {
	tkn, err := auth.GetToken(loginOpts.GitpodURL)
	if err != nil {
		return nil, err
	}

	wshost := loginOpts.GitpodURL
	wshost = strings.ReplaceAll(wshost, "https://", "wss://")
	wshost = strings.ReplaceAll(wshost, "http://", "ws://")
	wshost += "/api/v1"

	for {
		if tkn == "" {
			tkn, err = auth.Login(context.Background(), loginOpts)
			if tkn != "" {
				err = auth.SetToken(loginOpts.GitpodURL, tkn)
				if err != nil {
					logrus.WithField("origin", loginOpts.GitpodURL).Warnf("could not write token to keyring: %s", err)
					// Allow to continue
					err = nil
				}
			}
			if err != nil {
				return nil, err
			}
		}

		client, err := gitpod.ConnectToServer(wshost, gitpod.ConnectToServerOpts{
			Context:             context.Background(),
			Token:               tkn,
			Log:                 logrus.NewEntry(logrus.StandardLogger()),
			ReconnectionHandler: reconnectionHandler,
		})
		if err != nil {
			return nil, err
		}
		valid, validErr := validateToken(client)
		if valid {
			return client, nil
		}
		logrus.WithField("origin", loginOpts.GitpodURL).Error("gitpod token is invalid")
		err = client.Close()
		if err != nil {
			logrus.WithError(err).WithField("origin", loginOpts.GitpodURL).Warn("failed to close connection to gitpod server")
		}

		tkn = ""
		err = auth.DeleteToken(loginOpts.GitpodURL)
		if err != nil {
			logrus.WithError(err).WithField("origin", loginOpts.GitpodURL).Warn("failed to delete gitpod token")
		}

		if validErr != nil {
			return nil, validErr
		}
	}
}

func validateToken(client *gitpod.APIoverJSONRPC) (bool, error) {
	scopes, err := client.GetScopes(context.Background())
	if e, ok := err.(*gitpod.ErrBadHandshake); ok && e.Resp.StatusCode == 401 {
		return false, nil
	}
	if err != nil && strings.Contains(err.Error(), "jsonrpc2: code 403") {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	scopesMap := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		scopesMap[scope] = struct{}{}
	}
	for _, scope := range auth.Scopes {
		_, ok := scopesMap[scope]
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

type logCallbacks struct{}

func (*logCallbacks) InstanceUpdate(w *bastion.Workspace) {
	logrus.WithField("workspace", w).Info("instance update")
}
