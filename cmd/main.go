package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	port = ":9999"
)

type Client struct {
	*websocket.Conn
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	clients []*Client
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("websocket connection died")
			return
		}

		cl := &Client{conn}
		fmt.Printf("%s connected\n", conn.RemoteAddr())
		clients = append(clients, cl)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				log.Fatal("websocket reader error")
				return
			}
			fmt.Printf("%s send: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range clients {
				if err = client.WriteMessage(msgType, msg); err != nil {
					log.Fatal("websocket writer error")
					return
				}
			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	log.Printf("http server started on port %s", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("http server died")
	}
}
