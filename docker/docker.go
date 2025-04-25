package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type ContainerConfig struct {
	Image  string
	Cmd    []string
	Mounts []mount.Mount
	Name   string
}

type DockerContainerManager struct {
	client *client.Client
}

func NewDockerContainerManager() (*DockerContainerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerContainerManager{client: cli}, nil
}

func (m *DockerContainerManager) CreateAndStart(config ContainerConfig) (types.HijackedResponse, string, error) {
	ctx := context.Background()

	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image:        config.Image,
		Cmd:          config.Cmd,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts:     config.Mounts,
		AutoRemove: true,
	}, nil, nil, config.Name)
	if err != nil {
		return types.HijackedResponse{}, "", err
	}

	attach, err := m.client.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return types.HijackedResponse{}, "", err
	}

	err = m.client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return types.HijackedResponse{}, "", err
	}

	return attach, resp.ID, nil
}
