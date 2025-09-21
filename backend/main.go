package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	Clients  = make(map[*websocket.Conn]bool)
	Messages = make(chan *Message)
	Rooms    = make(map[string]Room)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()
	Clients[ws] = true
	log.Println("Client connected:", len(Clients))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(Clients, ws)
			break
		}

		event, err := NewEvent(msg)
		if err != nil {
			log.Println(err)
			break
		}
		HandleEvents(*event, ws)
	}
}

func handleMessages() {
	for {
		msg := <-Messages

		switch msg.Type {
		case GlobalBroadcast:
			for client := range Clients {
				err := client.WriteMessage(websocket.TextMessage, msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					client.Close()
					delete(Clients, client)
				}
			}
		case ClientSpecific:
			err := msg.Client.WriteMessage(websocket.TextMessage, msg.Content)
			if err != nil {
				log.Println("Write error:", err)
				msg.Client.Close()
				delete(Clients, msg.Client)
			}
		case RoomBroadcast:
			// TODO: Implement later
		}

		log.Println(string(msg.Content))
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
