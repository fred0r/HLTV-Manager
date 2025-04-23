package main

import (
	"HLTV-Manager/reader"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var testHLTV reader.HLTV

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("Docker client error:", err)
		return
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println("Getwd error:", err)
		return
	}

	demoPath := filepath.Join(path, testHLTV.Name, "cstrike")

	os.MkdirAll(demoPath, 0755)

	// Создание контейнера
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "my-hltv",
		Cmd:          []string{"+connect", testHLTV.Connect, "-port", testHLTV.HltvPort, "+record", testHLTV.DemoName},
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: demoPath,
				Target: "/home/hltv/cstrike",
			},
		},
		AutoRemove: true,
	}, nil, nil, "")
	if err != nil {
		log.Println("Create error:", err)
		return
	}

	// Подключаемся к stdin/stdout контейнера
	attach, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		log.Println("Attach error:", err)
		return
	}
	defer attach.Close()

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Println("Start error:", err)
		return
	}

	// Вывод контейнера -> WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := attach.Reader.Read(buf)
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				log.Println("Write WS error:", err)
				break
			}
		}
	}()

	// WebSocket -> stdin контейнера
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			break
		}
		_, err = attach.Conn.Write(msg)
		if err != nil {
			log.Println("Docker write error:", err)
			break
		}
	}
}

func exampleRunDocker() {
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
