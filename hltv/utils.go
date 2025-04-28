package hltv

import (
	"HLTV-Manager/config"
	log "HLTV-Manager/logger"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func createDemosDir(id int) (string, error) {
	demoPath := filepath.Join(config.HltvDemosDir(), "demos", strconv.Itoa(id), "cstrike")

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

func createHltvCfg(id int, cvars []string) (string, error) {
	cfgPath := filepath.Join(config.HltvDemosDir(), "demos", strconv.Itoa(id), "hltv.cfg")

	file, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при открытии hltv.cfg (%d): %v", id, err)
		return "", err
	}
	defer file.Close()

	fileContent := strings.Join(cvars, "\n") + "\n"
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

func parseDemoFilename(demoname string, filename string) (Demos, error) {
	name := strings.TrimSuffix(filename, ".dem")

	name = strings.TrimPrefix(name, fmt.Sprintf("%s-", demoname))

	parts := strings.SplitN(name, "-", 2)
	if len(parts) != 2 {
		return Demos{}, fmt.Errorf("неправильный формат имени файла")
	}

	datetime := parts[0]
	mapName := parts[1]

	if len(datetime) < 10 {
		return Demos{}, fmt.Errorf("неправильный формат даты/времени")
	}
	datePart := datetime[:6]
	timePart := datetime[6:]

	date := fmt.Sprintf("20%s.%s.%s", datePart[:2], datePart[2:4], datePart[4:6])
	time := fmt.Sprintf("%s:%s", timePart[:2], timePart[2:4])

	return Demos{
		Name: filename,
		Date: date,
		Time: time,
		Map:  mapName,
	}, nil
}
