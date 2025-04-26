package main

import (
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

	read, err := reader.ReadHLTVRunners("hltv-runners.yaml")
	if err != nil {
		log.ErrorLogger.Printf("Обшика в чтении данных конфига hltv: %v", err)
		return
	}

	log.InfoLogger.Println(read)

	hltv := &[]hltv.HLTV{{
		ID: 1,
		Settings: hltv.Settings{
			Name:     read[0].Name,
			Connect:  read[0].Connect,
			Port:     read[0].Port,
			DemoName: read[0].DemoName,
			Config: []string{
				"+exec hltv.cfg",
				"+sv_lan 0",
				"+maxclients 32",
			},
		},
		Demos: []hltv.Demos{
			{
				Map:  "de_dust2",
				Date: "2025.04.25",
				Time: "12:44",
			},
			{
				Map:  "de_dust",
				Date: "2025.04.25",
				Time: "12:55",
			},
		},
		Docker: nil,
	},
		{
			ID: 1,
			Settings: hltv.Settings{
				Name:     "Second HLTV",
				Connect:  "127.0.0.1:27016",
				Port:     "27021",
				DemoName: "second_demo",
				Config: []string{
					"+exec hltv.cfg",
					"+sv_lan 0",
					"+maxclients 32",
				},
			},
			Demos: []hltv.Demos{
				{
					Map:  "de_dust2",
					Date: "2025.04.25",
					Time: "12:44",
				},
				{
					Map:  "de_dust",
					Date: "2025.04.25",
					Time: "12:55",
				},
			},
			Docker: nil,
		},
	}

	// hltvConfig := hltv.Settings{
	// 	Name:     read[0].Name,
	// 	Connect:  read[0].Connect,
	// 	Port:     read[0].Port,
	// 	DemoName: read[0].DemoName,
	// 	Config:   read[0].Config,
	// }

	// hltv, err := hltv.NewHLTV(1, hltvConfig)
	// if err != nil {
	// 	log.ErrorLogger.Printf("Ошибка при создании HLTV: %v", err)
	// 	return
	// }

	// err = hltv.Start()
	// if err != nil {
	// 	log.ErrorLogger.Printf("Ошибка запуске HLTV: %v", err)
	// 	return
	// }

	// go hltv.ShowTerminal()

	site := &site.Site{HLTV: hltv}
	go site.Init()

	log.InfoLogger.Printf("Starting site")
	err = http.ListenAndServe("localhost:3002", nil)
	if err != nil {
		log.ErrorLogger.Fatal(err)
	}

	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutDown

		//hltv.Quit()

		fmt.Println("Программа завершена.")

		os.Exit(0)
	}()

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
