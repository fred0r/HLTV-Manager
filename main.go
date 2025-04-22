package main

import (
	log "HLTV-Manager/logger"
	"HLTV-Manager/reader"
	"fmt"
)

var testHLTV reader.HLTV

func main() {
	err := log.InitLogger("./log/")
	if err != nil {
		fmt.Println("Ошибка при инициализации логгера: ", err)
		return
	}

	HLTV, err := reader.GetHLTV("hltv-runners.yaml")
	if err != nil {
		log.WarningLogger.Printf("Обшика в чтении данных конфига hltv: %v", err)
		return
	}

	log.InfoLogger.Println(HLTV)

	testHLTV = HLTV[0]

	exampleRunDocker()
}
