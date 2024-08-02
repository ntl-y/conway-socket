package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	clients []*websocket.Conn
)

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
	http.HandleFunc("/game", websocketPage)
	http.HandleFunc("/", home)

	log.Println("server started at port :3000")
	http.ListenAndServe(":3000", nil)
}
