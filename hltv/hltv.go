package hltv

import (
	"HLTV-Manager/docker"
	log "HLTV-Manager/logger"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
)

const maxLogLines = 100

type HLTV struct {
	ID        int64
	Settings  Settings
	Log       []string
	mu        sync.Mutex
	Container *docker.DockerContainerManager
	Attach    types.HijackedResponse
	ContID    string
}

type Settings struct {
	Name     string
	Connect  string
	Port     string
	DemoName string
	Config   []string
}

func NewHLTV(id int64, settings Settings) (*HLTV, error) {
	containerManager, err := docker.NewDockerContainerManager()
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при инициализации контейнера hltv (%d): %v", id, err)
		return nil, err
	}

	return &HLTV{
		ID:        id,
		Settings:  settings,
		Container: containerManager,
	}, nil
}

func (hltv *HLTV) Start() error {
	path, err := os.Getwd()
	if err != nil {
		log.ErrorLogger.Printf("Обшика при получении пути hltv (%d): %v", hltv.ID, err)
		return err
	}

	demoPath, err := createDemosDir(path, hltv.ID)
	if err != nil {
		return err
	}

	cfgPath, err := createHltvCfg(path, hltv.ID, hltv.Settings.Config)
	if err != nil {
		return err
	}

	cmd := []string{
		"+connect", hltv.Settings.Connect,
		"-port", hltv.Settings.Port,
		"+record", hltv.Settings.DemoName,
	}

	hltv.Attach, hltv.ContID, err = hltv.Container.CreateAndStart(docker.ContainerConfig{
		Image: "ghcr.io/wesstorn/hltv-files:v1.0",
		Cmd:   cmd,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfgPath,
				Target: "/home/hltv/hltv.cfg",
			},
			{
				Type:   mount.TypeBind,
				Source: demoPath,
				Target: "/home/hltv/cstrike",
			},
		},
		Name: "hltv_" + strconv.FormatInt(hltv.ID, 10),
	})

	if err != nil {
		log.ErrorLogger.Printf("Обшика при запуске контейнера hltv (%d): %v", hltv.ID, err)
		return err
	}
	go func() { // READER CONTAINER TODO: Вынести отдельно и использовать по надобности
		buf := make([]byte, 1024)
		for {
			n, err := hltv.Attach.Reader.Read(buf)
			if err != nil {
				break
			}
			line := string(buf[:n])

			hltv.mu.Lock()
			hltv.Log = append(hltv.Log, line)
			if len(hltv.Log) > maxLogLines {
				hltv.Log = hltv.Log[len(hltv.Log)-maxLogLines:]
			}

			fmt.Println(line)
			hltv.mu.Unlock()
		}
	}()

	return nil
}

func (hltv *HLTV) Quit() error {
	err := hltv.WriteCommand("quit")
	if err != nil {
		return err
	}

	if closer, ok := hltv.Attach.Conn.(interface{ CloseWrite() error }); ok {
		_ = closer.CloseWrite()
	}

	hltv.Attach.Close()

	return nil
}

func (hltv *HLTV) GetLog() []string {
	hltv.mu.Lock()
	defer hltv.mu.Unlock()
	return append([]string{}, hltv.Log...) // копия
}

func (hltv *HLTV) WriteCommand(cmd string) error {
	_, err := hltv.Attach.Conn.Write([]byte(cmd + "\n"))
	return err
}
