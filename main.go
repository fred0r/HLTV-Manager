package main

import (
	"HLTV-Manager/config"
	"HLTV-Manager/hltv"
	log "HLTV-Manager/logger"
	"HLTV-Manager/reader"
	"HLTV-Manager/site"
	"fmt"
	"net/http"
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

	err = config.InitConfig()
	if err != nil {
		log.ErrorLogger.Printf("Ошибка при инициализации конфигурации: %v", err)
		return
	}

	read, err := reader.ReadHLTVRunners()
	if err != nil {
		log.ErrorLogger.Printf("Обшика в чтении данных конфига hltv: %v", err)
		return
	}

	log.InfoLogger.Println(read)

	var hltvs []*hltv.HLTV

	for i, runner := range read {
		hltvConfig := hltv.Settings{
			Name:       runner.Name,
			Connect:    runner.Connect,
			Port:       runner.Port,
			DemoName:   runner.DemoName,
			MaxDemoDay: runner.MaxDemoDay,
			Cvars:      runner.Cvars,
		}

		h, err := hltv.NewHLTV(i+1, hltvConfig)
		if err != nil {
			log.ErrorLogger.Printf("Ошибка при создании HLTV %d: %v", i, err)
			continue
		}

		err = h.Start()
		if err != nil {
			log.ErrorLogger.Printf("Ошибка при запуске HLTV %d: %v", i, err)
			continue
		}

		hltvs = append(hltvs, h)

		go h.ShowTerminal()
	}

	site := &site.Site{HLTV: hltvs}
	go site.Init()

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutDown

		for _, hltv := range hltvs {
			hltv.Quit()
		}

		fmt.Println("Программа завершена.")

		os.Exit(0)
	}()

	address := fmt.Sprintf("%s:%s", config.SiteIP(), config.SitePort())
	log.InfoLogger.Println("Starting site: ", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.ErrorLogger.Fatal(err)
	}

	// for {
	// 	in := bufio.NewReader(os.Stdin)
	// 	line, err := in.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("ERR")
	// 		continue
	// 	}

	// 	err = hltv.WriteCommand(line)
	// 	if err != nil {
	// 		fmt.Println("Write to container error:", err)
	// 		break
	// 	}
	// }

	select {}
}
