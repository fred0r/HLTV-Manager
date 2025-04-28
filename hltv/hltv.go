package hltv

import (
	"HLTV-Manager/docker"
	log "HLTV-Manager/logger"
	"fmt"
	"strings"
)

const maxLogLines = 100

type HLTV struct {
	ID       int
	Settings Settings
	Demos    []Demos
	Docker   *docker.Docker
}

type Settings struct {
	Name       string
	Connect    string
	Port       string
	DemoDir    string
	DemoName   string
	MaxDemoDay string
	Cvars      []string
}

type Demos struct {
	ID   int
	Name string
	Date string
	Time string
	Map  string
}

func NewHLTV(id int, settings Settings) (*HLTV, error) {
	docker, err := docker.NewDockerClient()
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error creating Docker client: %v", id, settings.Name, err)
		return nil, err
	}

	return &HLTV{
		ID:       id,
		Settings: settings,
		Docker:   docker,
	}, nil
}

func (hltv *HLTV) Start() error {
	var err error

	hltv.Settings.DemoDir, err = createDemosDir(hltv)
	if err != nil {
		return err
	}

	cfgPath, err := createHltvCfg(hltv)
	if err != nil {
		return err
	}

	hltvData := docker.Hltv{
		ID:   hltv.ID,
		Name: hltv.Settings.Name,
	}

	err = hltv.Docker.CreateAndStart(docker.HltvContainerConfig{
		Cmd: []string{
			"+connect", hltv.Settings.Connect,
			"-port", hltv.Settings.Port,
			"+record", hltv.Settings.DemoName,
		},
		DemoPath: hltv.Settings.DemoDir,
		CfgPath:  cfgPath,
		Hltv:     hltvData,
	})
	if err != nil {
		return err
	}

	err = hltv.DemoControl()
	if err != nil {
		return err
	}

	return nil
}

func (hltv *HLTV) Quit() error {
	err := hltv.WriteCommand("quit")
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Failed to write quit command: %v", hltv.ID, hltv.Settings.Name, err)
		return err
	}

	if closer, ok := hltv.Docker.Attach.Conn.(interface{ CloseWrite() error }); ok {
		_ = closer.CloseWrite()
	}

	hltv.Docker.Attach.Close()

	return nil
}

func (hltv *HLTV) ShowTerminal() {
	buf := make([]byte, 1024)
	for {
		n, err := hltv.Docker.Attach.Reader.Read(buf)
		if err != nil {
			break
		}
		line := string(buf[:n])
		line = strings.TrimRight(line, "\n")
		fmt.Println(line)
	}
}

func (hltv *HLTV) WriteCommand(cmd string) error {
	_, err := hltv.Docker.Attach.Conn.Write([]byte(cmd + "\n"))
	return err
}
