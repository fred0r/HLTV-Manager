package main

import (
	"HLTV-Manager/hltv"
	log "HLTV-Manager/logger"
	"HLTV-Manager/reader"
	"bufio"
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
		log.ErrorLogger.Printf("Обшика в чтении данных конфига hltv: %v", err)
		return
	}

	log.InfoLogger.Println(read)

	hltvConfig := hltv.Config{
		Connect:  read[0].Connect,
		HltvPort: read[0].HltvPort,
		DemoFile: read[0].Name,
		DemoName: read[0].DemoName,
	}

	hltv, err := hltv.NewHLTV(1, hltvConfig)
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при создании HLTV: %v", err)
		return
	}

	err = hltv.Start()
	if err != nil {
		log.ErrorLogger.Printf("Ошибка запуске HLTV: %v", err)
		return
	}

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutDown

		hltv.Quit()

		fmt.Println("Программа завершена.")

		os.Exit(0)
	}()

	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			fmt.Println("ERR")
			continue
		}

		err = hltv.WriteCommand(line)
		if err != nil {
			fmt.Println("Write to container error:", err)
			break
		}
	}

	select {}
}
