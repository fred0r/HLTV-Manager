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
		http.Error(w, "Invalid url.", http.StatusBadRequest)
		return
	}

	hltvID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid HLTV ID.", http.StatusBadRequest)
		return
	}

	demoID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid Demo ID.", http.StatusBadRequest)
		return
	}

	var hltv *hltv.HLTV
	for _, h := range site.HLTV {
		if h.ID == hltvID {
			hltv = h
			break
		}
	}

	if hltv == nil {
		http.Error(w, "HLTV server not found.", http.StatusNotFound)
		return
	}

	demoName, err := hltv.GetDemoFileName(demoID)
	if err != nil {
		http.Error(w, "Invalid GetDemoFileName Demo.", http.StatusBadRequest)
		return
	}

	demoFilePath := filepath.Join(hltv.Settings.DemoDir, demoName)
	if _, err := os.Stat(demoFilePath); os.IsNotExist(err) {
		http.Error(w, "Demo file not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+demoName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, demoFilePath)
}
