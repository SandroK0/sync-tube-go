package api

import (
	"log"
	"net/http"

	"github.com/SandroK0/sync-tube-go/backend/entities"
	"github.com/gorilla/websocket"
)

var (
	Clients  = make(map[*websocket.Conn]bool)
	Messages = make(chan *entities.Message)
	Rooms    = make(map[string]*entities.Room)
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

		event, err := NewEvent(msg)
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
		case entities.GlobalBroadcast:
			for client := range Clients {
				err := client.WriteMessage(websocket.TextMessage, msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					client.Close()
					delete(Clients, client)
				}
			}
		case entities.ClientSpecific:
			err := msg.Client.WriteMessage(websocket.TextMessage, msg.Content)
			if err != nil {
				log.Println("Write error:", err)
				msg.Client.Close()
				delete(Clients, msg.Client)
			}
		case entities.RoomBroadcast:
			room, exists := Rooms[msg.RoomName]
			if !exists {
				log.Println("Room not found:", msg.RoomName)
				return
			}

			for _, user := range room.Users {
				err := user.Conn.WriteMessage(websocket.TextMessage, msg.Content)
				if err != nil {
					log.Println("Write error:", err)
					user.Conn.Close()
					delete(Clients, user.Conn)
				}
			}
		}

		log.Println(string(msg.Content))
	}
}
