package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

type Event struct {
	EventType string `json:"eventType"`
	Data      any    `json:"data"`
}

func newEvent(msg []byte) (*Event, error) {

	var event Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &event, nil
}

func handleEvents(event Event, ws *websocket.Conn) {
	if event.EventType == "test" {
		broadcast <- `{"eventType": "test"}`
	}

	switch event.EventType {
	case "test":
		broadcast <- `{"eventType": "test"}`
	case "foo":
		ws.WriteMessage(websocket.TextMessage, []byte(`{"eventType": "foo only for u"}`))
	default:
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true
	log.Println("Client connected:", len(clients))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, ws)
			break
		}

		event, err := newEvent(msg)
		if err != nil {
			log.Println(err)
			break
		}

		handleEvents(*event, ws)

	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
