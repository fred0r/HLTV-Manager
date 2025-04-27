package site

import (
	"HLTV-Manager/hltv"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (site *Site) downloadHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Неверный путь", http.StatusBadRequest)
		return
	}

	hltvID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Неверный ID HLTV", http.StatusBadRequest)
		return
	}

	demoID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Неверный ID Demo", http.StatusBadRequest)
		return
	}

	var hltv *hltv.HLTV
	for _, h := range site.HLTV {
		if h.ID == int64(hltvID) {
			hltv = h
			break
		}
	}

	if hltv == nil {
		http.Error(w, "HLTV сервер не найден", http.StatusNotFound)
		return
	}

	demoName, err := hltv.GetDemoFileName(demoID)
	if err != nil {
		http.Error(w, "Неверный GetDemoFileName Demo", http.StatusBadRequest)
		return
	}

	demoFilePath := filepath.Join(hltv.Settings.DemoDir, demoName)
	if _, err := os.Stat(demoFilePath); os.IsNotExist(err) {
		http.Error(w, "Файл демки не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+demoName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, demoFilePath)
}
