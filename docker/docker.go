package docker

import (
	"context"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type HltvContainerConfig struct {
	Cmd      []string
	DemoPath string
	CfgPath  string
	Mounts   []mount.Mount
	HltvID   int64
}

type Docker struct {
	client *client.Client
	Attach types.HijackedResponse
}

func NewDockerClient() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{client: cli}, nil
}

func (docker *Docker) CreateAndStart(config HltvContainerConfig) error {
	ctx := context.Background()

	resp, err := docker.client.ContainerCreate(ctx, &container.Config{
		Image:        "ghcr.io/wesstorn/hltv-files:v1.0", // TODO: Add config
		Cmd:          config.Cmd,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: config.DemoPath,
				Target: "/home/hltv/cstrike",
			},
			{
				Type:   mount.TypeBind,
				Source: config.CfgPath,
				Target: "/home/hltv/hltv.cfg",
			},
		},
		AutoRemove: true,
	}, nil, nil, "hltv_"+strconv.FormatInt(config.HltvID, 10))
	if err != nil {
		return err
	}

	docker.Attach, err = docker.client.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return err
	}

	err = docker.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}
