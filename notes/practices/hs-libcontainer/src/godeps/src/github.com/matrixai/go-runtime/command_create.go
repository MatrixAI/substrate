package main

import (
	"github.com/coreos/go-systemd/activation"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/specconv"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path/filepath"
)

// CreateCommand handles the creation of containers
type createCommand struct {
	runnableCommand

	bundle        string
	consoleSocket string
	pidFile       string
	preserveFds   int
}

// Execute sets up the environment for the container.
func (cmd *createCommand) Execute() (interface{}, error) {
	// convert pid-file to an absolute path so we can write to
	// the right file after chdir to bundle
	pidFile, err := filepath.Abs(cmd.pidFile)
	if err != nil {
		return nil, err
	}
	cmd.pidFile = pidFile
	spec, err := setupSpec(cmd.bundle)
	if err != nil {
		return nil, err
	}
	status, err := cmd.startContainer(spec)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (cmd *createCommand) createContainer(spec *specs.Spec) (libcontainer.Container, error) {
	rootless, err := isRootless(cmd.rootless)
	if err != nil {
		return nil, err
	}
	config, err := specconv.CreateLibcontainerConfig(
		&specconv.CreateOpts{
			CgroupName:       cmd.id,
			UseSystemdCgroup: cmd.systemdCgroup,
			NoPivotRoot:      cmd.noPivot,
			NoNewKeyring:     cmd.noNewKeyring,
			Spec:             spec,
			Rootless:         rootless,
		})
	if err != nil {
		return nil, err
	}
	factory, err := cmd.loadFactory()
	if err != nil {
		return nil, err
	}
	return factory.Create(cmd.id, config)
}

func (cmd *createCommand) startContainer(spec *specs.Spec) (int, error) {
	notifySocket := newNotifySocket(cmd.statePath, cmd.notifySocket, cmd.id)
	if notifySocket != nil {
		notifySocket.setupSpec(spec)
	}
	container, err := cmd.createContainer(spec)
	if err != nil {
		return -1, err
	}
	if notifySocket != nil {
		if err := notifySocket.setupSocket(); err != nil {
			return -1, err
		}
	}
	// Support on-demand socket activation by passing
	// file descriptors into the container init process
	// TODO: Change this into something that doesn't use env var
	listenFDs := []*os.File{}
	if os.Getenv("LISTEN_FDS") != "" {
		listenFDs = activation.Files(false)
	}
	r := &runner{
		enableSubreaper: true, // Set current process as subreaper
		shouldDestroy:   true,
		container:       container,
		listenFDs:       listenFDs,
		notifySocket:    notifySocket,
		consoleSocket:   cmd.consoleSocket,
		detach:          false,
		pidFile:         cmd.pidFile,
		preserveFDs:     cmd.preserveFds,
		action:          CT_ACT_CREATE,
		criuOpts:        nil,
		init:            true,
	}
	return r.run(spec.Process)
}
