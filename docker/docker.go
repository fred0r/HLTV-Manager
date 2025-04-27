package docker

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

	err := docker.StopContainerIfExists(config.HltvID)
	if err != nil {
		return err
	}

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

func (docker *Docker) StopContainerIfExists(hltvID int64) error {
	ctx := context.Background()

	containerName := "hltv_" + strconv.FormatInt(hltvID, 10)
	containers, err := docker.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	fmt.Println("Start stop and removed existing container:", containerName)

	for _, cont := range containers {
		if cont.Names[0] == "/"+containerName {

			err := docker.client.ContainerStop(ctx, cont.ID, container.StopOptions{})
			if err != nil {
				return fmt.Errorf("failed to stop container: %w", err)
			}

			tries := 5
			for i := 0; i < tries; i++ {
				_, err := docker.client.ContainerInspect(ctx, cont.ID)
				if err != nil {
					if client.IsErrNotFound(err) {
						fmt.Println("Container already removed, skipping remove.")
						break
					}
					return fmt.Errorf("failed to inspect container: %w", err)
				}

				err = docker.client.ContainerRemove(ctx, cont.ID, container.RemoveOptions{Force: true})
				if err == nil {
					fmt.Println("Removed container:", containerName)
					break
				} else {
					fmt.Printf("Error removing container, retrying in 3 seconds... (Attempt %d/%d)\n", i+1, tries)
					time.Sleep(3 * time.Second)
				}
			}

			if err != nil {
				return fmt.Errorf("failed to remove container after %d attempts: %w", tries, err)
			}
		}
	}

	return nil
}
