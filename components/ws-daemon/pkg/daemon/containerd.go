// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package daemon

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/xerrors"
	"k8s.io/client-go/kubernetes"

	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/gitpod-io/gitpod/ws-daemon/pkg/container"
	"github.com/gitpod-io/gitpod/ws-daemon/pkg/dispatch"
	"github.com/prometheus/procfs"
)

const (
	maxDeletionAttempts    = 10
	unmountAttemptInterval = 2 * time.Second
)

// ContainerdUnmount ensures the mark mount is removed
type ContainerdUnmount struct {
	mu      sync.Mutex
	handled map[string]struct{}
}

// WorkspaceAdded does nothing but implemented the dispatch.Listener interface
func (c *ContainerdUnmount) WorkspaceAdded(ctx context.Context, ws *dispatch.Workspace) error {
	return nil
}

// WorkspaceUpdated gets called when a workspace pod is updated. For containers being deleted, we'll check
// if they're still running after their terminationGracePeriod and if Kubernetes still knows about them.
func (c *ContainerdUnmount) WorkspaceUpdated(ctx context.Context, ws *dispatch.Workspace) error {
	if ws.Pod.DeletionTimestamp == nil {
		return nil
	}

	c.mu.Lock()
	if c.handled == nil {
		c.handled = make(map[string]struct{})
	}
	if _, exists := c.handled[ws.InstanceID]; exists {
		c.mu.Unlock()
		return nil
	}
	c.handled[ws.InstanceID] = struct{}{}
	c.mu.Unlock()

	dsp := dispatch.GetFromContext(ctx)
	go func() {
		err := c.ensureMarkUnmount(dsp.Runtime, dsp.Kubernetes, ws)
		if err != nil {
			log.WithError(err).Error("cannot unmount mark")
		}
	}()

	return nil
}

// ensurePodGetsDeleted will check if the container still exists on this node, i.e. still runs.
// If it doesn't, it'll force unmount of the mark mount.
func (c *ContainerdUnmount) ensureMarkUnmount(rt container.Runtime, clientSet kubernetes.Interface, ws *dispatch.Workspace) (err error) {
	var (
		log         = log.WithFields(ws.OWI())
		containerID = ws.ContainerID
	)

	delay := unmountAttemptInterval
	for attempt := 0; attempt < maxDeletionAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(delay)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		var exists bool
		exists, err = rt.ContainerExists(ctx, containerID)
		cancel()
		if err != nil {
			log.WithField("attempt", attempt).WithError(err).Warn("cannot check if container still exists")
			continue
		}
		if exists {
			continue
		}

		if err := unmountMark(ws.WorkspaceID); err != nil {
			log.WithError(err).WithField("workspaceId", ws.WorkspaceID).Error("cannot unmount mark mount")
			return err
		}

		return nil
	}

	return err
}

// if the mark mount still exists in /proc/mounts it means we failed to unmount it and
// we cannot remove the content. As a side effect the pod will stay in Terminating state
func unmountMark(path string) error {
	mounts, err := procfs.GetMounts()
	if err != nil {
		return xerrors.Errorf("unexpected error reading /proc/mounts: %w", err)
	}

	for _, mount := range mounts {
		if !strings.Contains(mount.MountPoint, path) {
			continue
		}

		log.WithField("path", path).Debug("Unmounting pending mark")
		err = unix.Unmount(mount.MountPoint, 0)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
