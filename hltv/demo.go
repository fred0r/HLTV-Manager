package hltv

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (h *HLTV) DemoControl() error {
	err := h.LoadDemosFromFolder()
	if err != nil {
		return fmt.Errorf("ошибка обхода папки: %w", err)
	}

	sort.Slice(h.Demos, func(i, j int) bool {
		dateI, _ := time.Parse("2006.01.02 15:04", h.Demos[i].Date+" "+h.Demos[i].Time)
		dateJ, _ := time.Parse("2006.01.02 15:04", h.Demos[j].Date+" "+h.Demos[j].Time)
		return dateI.After(dateJ)
	})

	err = h.DeleteOldDemos()
	if err != nil {
		return fmt.Errorf("ошибка обхода папки: %w", err)
	}

	return nil
}

func (h *HLTV) LoadDemosFromFolder() error {
	var demos []Demos

	var id int

	err := filepath.Walk(h.Settings.DemoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".dem") {
			id++
			demo, err := parseDemoFilename(h.Settings.DemoName, info.Name())
			demo.ID = id
			if err != nil {
				fmt.Println("Ошибка парсинга файла:", info.Name(), err)
				return nil
			}
			demos = append(demos, demo)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка обхода папки: %w", err)
	}

	h.Demos = demos
	return nil
}

func (h *HLTV) DeleteOldDemos() error {
	now := time.Now()

	for _, demo := range h.Demos {
		demoDate, err := time.Parse("2006.01.02 15:04", demo.Date+" "+demo.Time)
		if err != nil {
			return fmt.Errorf("не удалось распарсить дату для демки %s: %w", demo.Name, err)
		}

		maxDemoDay, err := strconv.Atoi(h.Settings.MaxDemoDay)
		if err != nil {
			return fmt.Errorf("Ошибка конвертации %s: %w", demo.Name, err)
		}

		if now.Sub(demoDate).Hours() > float64(maxDemoDay*24) {
			demoPath := filepath.Join(h.Settings.DemoDir, demo.Name)
			err := os.Remove(demoPath)
			if err != nil {
				return fmt.Errorf("не удалось удалить демку %s: %w", demo.Name, err)
			}
			fmt.Printf("Удалена старая демка: %s\n", demo.Name)
		}
	}

	return nil
}

func (h *HLTV) GetDemoFileName(demoID int) (string, error) {
	var demo Demos
	for _, d := range h.Demos {
		if d.ID == demoID {
			demo = d
			break
		}
	}

	datePart := strings.ReplaceAll(demo.Date, ".", "")[2:]
	timePart := strings.ReplaceAll(demo.Time, ":", "")

	return fmt.Sprintf("%s-%s-%s.dem", h.Settings.DemoName, datePart+timePart, demo.Map), nil
}
