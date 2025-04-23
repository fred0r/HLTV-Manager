package hltv

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

const maxLogLines = 100

type HLTV struct {
	ID     int64
	Config Config
	Log    []string
	mu     sync.Mutex
	Attach types.HijackedResponse
}

type Config struct {
	Connect  string
	HltvPort string
	DemoFile string
	DemoName string
}

func (h *HLTV) Start(shutDown chan os.Signal) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println("Getwd error:", err)
		return err
	}

	demoPath := filepath.Join(path, h.Config.DemoFile, "cstrike")

	os.MkdirAll(demoPath, 0755)

	cmd := []string{
		"+connect", h.Config.Connect,
		"-port", h.Config.HltvPort,
		"+record", h.Config.DemoFile,
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "my-hltv",
		Cmd:          cmd,
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: demoPath,
				Target: "/home/hltv/cstrike",
			},
		},
		AutoRemove: true,
	}, nil, nil, "hltv_1") // TODO: Айди хлтв
	if err != nil {
		log.Println("Create error:", err)
		return err
	}

	h.Attach, err = cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		log.Println("Attach error:", err)
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Println("Start error:", err)
		return err
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := h.Attach.Reader.Read(buf)
			if err != nil {
				break
			}
			line := string(buf[:n])

			h.mu.Lock()
			h.Log = append(h.Log, line)
			if len(h.Log) > maxLogLines {
				h.Log = h.Log[len(h.Log)-maxLogLines:]
			}

			fmt.Println(line)
			h.mu.Unlock()
		}
	}()

	// go func() {
	// 	<-shutDown
	// 	err = h.WriteCommand("quit")
	// 	if err != nil {
	// 		fmt.Println("Write to container error:", err)
	// 		return
	// 	}
	// 	time.Sleep(5 * time.Second) // TODO: Костыль
	// 	h.Attach.Close()
	// 	fmt.Println("Программа завершена.")
	// 	os.Exit(0)
	// }()

	go func() {
		for {
			in := bufio.NewReader(os.Stdin)
			line, err := in.ReadString('\n')
			if err != nil {
				fmt.Println("ERR")
				continue
			}

			err = h.WriteCommand(line)
			if err != nil {
				fmt.Println("Write to container error:", err)
				break
			}
		}
	}()

	return nil
}

// func (h *HLTV) Quit() error {
// 	_, err := h.Attach.Conn.Write([]byte("quit"))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (h *HLTV) GetLog() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]string{}, h.Log...) // копия
}

func (h *HLTV) WriteCommand(cmd string) error {
	_, err := h.Attach.Conn.Write([]byte(cmd))
	return err
}
