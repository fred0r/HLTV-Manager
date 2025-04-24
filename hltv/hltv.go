package hltv

import (
	"HLTV-Manager/docker"
	log "HLTV-Manager/logger"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/docker/docker/api/types"
)

const maxLogLines = 100

type HLTV struct {
	ID        int64
	Config    Config
	Log       []string
	mu        sync.Mutex
	Container *docker.DockerContainerManager
	Attach    types.HijackedResponse
	ContID    string
}

type Config struct {
	Connect  string
	HltvPort string
	DemoFile string
	DemoName string
}

func NewHLTV(id int64, config Config) (*HLTV, error) {
	containerManager, err := docker.NewDockerContainerManager()
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при инициализации контейнера hltv (%d): %v", id, err)
		return nil, err
	}

	return &HLTV{
		ID:        id,
		Config:    config,
		Container: containerManager,
	}, nil
}

func (h *HLTV) Start() error {
	path, err := os.Getwd()
	if err != nil {
		log.ErrorLogger.Printf("Обшика при получении пути hltv (%d): %v", h.ID, err)
		return err
	}

	demoPath := filepath.Join(path, h.Config.DemoFile, "cstrike")

	err = os.MkdirAll(demoPath, 0755)
	if err != nil {
		log.ErrorLogger.Printf("Обшика при создании директории для демок hltv (%d): %v", h.ID, err)
		return err
	}

	err = os.Chown(demoPath, 1000, 1000)
	if err != nil {
		log.ErrorLogger.Printf("Обшика при выдаче прав директории для демок hltv (%d): %v", h.ID, err)
		return err
	}

	cmd := []string{
		"+connect", h.Config.Connect,
		"-port", h.Config.HltvPort,
		"+record", h.Config.DemoFile,
	}

	h.Attach, h.ContID, err = h.Container.CreateAndStart(docker.ContainerConfig{
		Image:     "my-hltv",
		Cmd:       cmd,
		MountHost: demoPath,
		MountCont: "/home/hltv/cstrike",
		Name:      "hltv_" + strconv.FormatInt(h.ID, 10),
	})

	if err != nil {
		log.ErrorLogger.Printf("Обшика при запуске контейнера hltv (%d): %v", h.ID, err)
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

	return nil
}

func (h *HLTV) Quit() error {
	err := h.WriteCommand("quit")
	if err != nil {
		return err
	}

	if closer, ok := h.Attach.Conn.(interface{ CloseWrite() error }); ok {
		_ = closer.CloseWrite()
	}

	h.Attach.Close()

	return nil
}

func (h *HLTV) GetLog() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]string{}, h.Log...) // копия
}

func (h *HLTV) WriteCommand(cmd string) error {
	_, err := h.Attach.Conn.Write([]byte(cmd + "\n"))
	return err
}
