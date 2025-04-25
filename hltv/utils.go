package hltv

import (
	log "HLTV-Manager/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func createDemosDir(path string, id int64) (string, error) {
	demoPath := filepath.Join(path, "demos", strconv.FormatInt(id, 10), "cstrike")

	err := os.MkdirAll(demoPath, 0755)
	if err != nil {
		log.ErrorLogger.Printf("Обшика при создании директории для демок hltv (%d): %v", id, err)
		return "", err
	}

	err = os.Chown(demoPath, 1000, 1000)
	if err != nil {
		log.ErrorLogger.Printf("Обшика при выдаче прав директории для демок hltv (%d): %v", id, err)
		return "", err
	}

	return demoPath, nil
}

func createHltvCfg(path string, id int64, config []string) (string, error) {
	cfgPath := filepath.Join(path, "demos", strconv.FormatInt(id, 10), "hltv.cfg")

	file, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при открытии hltv.cfg (%d): %v", id, err)
		return "", err
	}
	defer file.Close()

	fileContent := strings.Join(config, "\n") + "\n"
	_, err = file.Write([]byte(fileContent))
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при записи в hltv.cfg (%d): %v", id, err)
		return "", err
	}

	err = os.Chown(cfgPath, 1000, 1000)
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при выдаче прав файлу hltv.cfg (%d): %v", id, err)
		return "", err
	}

	return cfgPath, nil
}
