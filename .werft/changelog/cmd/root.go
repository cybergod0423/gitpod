package cmd

import (
	"fmt"

	"github.com/gitpod-io/gitpod/changelog/pkg/changelog"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var program = &cobra.Command{
	Use:   "changelog",
	Long:  "Little configurable CLI to generate the markdown for your changelogs from release-note blocks found into your project pull requests.",
	Short: "Generate markdown for your changelogs from release-note blocks.",
	Run: func(c *cobra.Command, args []string) {
		client := changelog.NewClient(opts.Token)
		notes, err := client.Get(opts.Org, opts.Repo, opts.Branch, opts.Milestone)
		if err != nil {
			logger.WithError(err).Fatal("error retrieving PRs")
		}
		output, err := changelog.Print(opts.Milestone, notes)
		if err != nil {
			logger.WithError(err).Fatal("error printing out release notes")
		}
		fmt.Println(output)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Setup flags before the command is initialized
	flags := program.PersistentFlags()
	flags.StringVarP(&opts.Milestone, "file", "f", opts.Milestone, "the existing changelog file")
	flags.StringVarP(&opts.Token, "token", "t", opts.Token, "a GitHub personal API token to perform authenticated requests")
}

func initConfig() {
	// nop
}

// Run ...
func Execute() error {
	return program.Execute()
}
