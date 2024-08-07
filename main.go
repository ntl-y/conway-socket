package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

var (
	clients []*websocket.Conn
)

type World struct {
	width  int
	height int
	area   []bool
}

func NewWorld(width int, height int) *World {
	return &World{
		width:  width,
		height: height,
		area:   make([]bool, width*height),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func listen(conn *websocket.Conn) {
	for {
		messageFromBufferType, messageFromBuffer, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%s send: %s", conn.RemoteAddr(), messageFromBuffer)

		for _, client := range clients {
			if err := client.WriteMessage(messageFromBufferType, messageFromBuffer); err != nil {
				client.Close()
				log.Println(err)
				return
			}
		}
	}
}
func websocketPage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clients = append(clients, conn)
	log.Println(conn.RemoteAddr(), " Connected!")

	listen(conn)
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	// http.HandleFunc("/game", websocketPage)
	// http.HandleFunc("/", home)

	cmd := exec.Command("sh", "-c", "env GOOS=js GOARCH=wasm go build -o static/game.wasm ./game")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Println("server started at port :8080")
	http.ListenAndServe(":8080", nil)
}
