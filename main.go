package main

import (
	"HLTV-Manager/hltv"
	log "HLTV-Manager/logger"
	"HLTV-Manager/reader"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := log.InitLogger("./log/")
	if err != nil {
		fmt.Println("Ошибка при инициализации логгера: ", err)
		return
	}

	read, err := reader.ReadHLTVRunners("hltv-runners.yaml")
	if err != nil {
		log.WarningLogger.Printf("Обшика в чтении данных конфига hltv: %v", err)
		return
	}

	log.InfoLogger.Println(read)

	hltv := &hltv.HLTV{
		ID: 1,
		Config: hltv.Config{
			Connect:  read[0].Connect,
			HltvPort: read[0].HltvPort,
			DemoFile: read[0].Name,
			DemoName: read[0].DemoName,
		},
	}

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGINT, syscall.SIGTERM)

	go hltv.Start(shutDown)

	select {}
}
