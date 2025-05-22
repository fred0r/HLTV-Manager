package reader

type Config struct {
	HLTV []HLTV `yaml:"HLTV"`
}

type HLTV struct {
	Name             string   `yaml:"Name"`
	ShowIP           string   `yaml:"ShowIP"`
	Connect          string   `yaml:"Connect"`
	Port             string   `yaml:"HltvPort"`
	GameID           string   `yaml:"GameID"`
	DemoName         string   `yaml:"DemoName"`
	MaxDemoDay       string   `yaml:"MaxDemoDay"`
	DebugTerminalLog bool     `yaml:"DebugTerminalLog"`
	Cvars            []string `yaml:"Cvars"`
}
