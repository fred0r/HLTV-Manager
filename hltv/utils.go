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

func createDemosDir(hltv *HLTV) (string, error) {
	demoPath := filepath.Join(config.HltvDemosDir(), "demos", strconv.Itoa(hltv.ID), "cstrike")

	err := os.MkdirAll(demoPath, 0755)
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error when creating a directory for demos: %v", hltv.ID, hltv.Settings.Name, err)
		return "", err
	}

	err = os.Chown(demoPath, 1000, 1000)
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error chown directory for demo: %v", hltv.ID, hltv.Settings.Name, err)
		return "", err
	}

	return demoPath, nil
}

func createHltvCfg(hltv *HLTV) (string, error) {
	cfgPath := filepath.Join(config.HltvDemosDir(), "demos", strconv.Itoa(hltv.ID), "hltv.cfg")

	file, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error when opening hltv.cfg: %v", hltv.ID, hltv.Settings.Name, err)
		return "", err
	}
	defer file.Close()

	fileContent := strings.Join(hltv.Settings.Cvars, "\n") + "\n"
	_, err = file.Write([]byte(fileContent))
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error when writing to hltv.cfg: %v", hltv.ID, hltv.Settings.Name, err)
		return "", err
	}

	err = os.Chown(cfgPath, 1000, 1000)
	if err != nil {
		log.ErrorLogger.Printf("HLTV (ID: %d, Name: %s) Error chown file for hltv.cfg: %v", hltv.ID, hltv.Settings.Name, err)
		return "", err
	}

	return cfgPath, nil
}

func parseDemoFilename(demoname string, filename string) (Demos, error) {
	name := strings.TrimSuffix(filename, ".dem")

	name = strings.TrimPrefix(name, fmt.Sprintf("%s-", demoname))

	parts := strings.SplitN(name, "-", 2)
	if len(parts) != 2 {
		return Demos{}, fmt.Errorf("incorrect file name format")
	}

	datetime := parts[0]
	mapName := parts[1]

	if len(datetime) < 10 {
		return Demos{}, fmt.Errorf("incorrect date/time format")
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
