package reader

type Config struct {
	HLTV []HLTV `yaml:"HLTV"`
}

type HLTV struct {
	Name       string   `yaml:"Name"`
	Connect    string   `yaml:"Connect"`
	Port       string   `yaml:"HltvPort"`
	DemoName   string   `yaml:"DemoName"`
	MaxDemoDay string   `yaml:"MaxDemoDay"`
	Cvars      []string `yaml:"Cvars"`
}
