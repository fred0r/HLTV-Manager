package reader

import (
	"io"
	"os"

	"HLTV-Manager/config"
	log "HLTV-Manager/logger"

	"gopkg.in/yaml.v3"
)

func ReadHLTVRunners() ([]HLTV, error) {
	file, err := os.OpenFile(config.HltvRunnerFile(), os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		if !os.IsNotExist(err) {
			log.ErrorLogger.Println("Error opening HLTV runners file:", err)
			return nil, err
		}
		log.ErrorLogger.Println("HLTV runners file does not exist:", err)
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.ErrorLogger.Println("Error reading HLTV runners file:", err)
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		log.ErrorLogger.Println("Error unmarshalling HLTV runners YAML:", err)
		return nil, err
	}

	// TODO: Debug

	return config.HLTV, nil
}
