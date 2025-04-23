package reader

import (
	"io"
	"os"

	log "HLTV-Manager/logger"

	"gopkg.in/yaml.v3"
)

func ReadHLTVRunners(filePath string) ([]HLTV, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		if !os.IsNotExist(err) {
			log.ErrorLogger.Println(err)
			return nil, err
		}
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}

	return config.HLTV, nil
}
