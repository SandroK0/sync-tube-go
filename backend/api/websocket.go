package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
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

		event, err := NewClientEvent(msg)
		if err != nil {
			log.Println(err)
			break
		}
		HandleEvents(*event, ws)
	}
}

func HandleMessages() {
	for {
		msg := <-Messages

		switch msg.Type {
		case GlobalBroadcast:
			for client := range Clients {
				err := client.WriteJSON(msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					client.Close()
					delete(Clients, client)
				}
			}
		case ClientSpecific:
			err := msg.Client.WriteJSON(msg.Content)
			if err != nil {
				log.Println("Write error:", err)
				msg.Client.Close()
				delete(Clients, msg.Client)
			}
		case RoomBroadcast:
			room, exists := Rooms[msg.RoomName]
			if !exists {
				log.Println("Room not found:", msg.RoomName)
				return
			}

			for _, user := range room.Users {
				err := user.Conn.WriteJSON(msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					user.Conn.Close()
					delete(Clients, user.Conn)
				}
			}
		}

	}
}
