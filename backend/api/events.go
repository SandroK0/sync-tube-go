package api

import (
	"encoding/json"
	"log"

	"github.com/SandroK0/sync-tube-go/backend/entities"
	"github.com/gorilla/websocket"
)

type Event struct {
	EventType string          `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

type JoinEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

func NewEvent(msg []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func HandleEvents(event Event, ws *websocket.Conn) {
	switch event.EventType {
	case "join":

		var joinData JoinEventData
		if err := json.Unmarshal(event.Data, &joinData); err != nil {
			log.Println("Error unmarshaling join event:", err)
			return
		}

		if joinData.RoomName == "" || joinData.Username == "" {
			log.Println("Missing roomName or username in join event")
			return
		}

		user := entities.NewUser(joinData.Username, ws)
		room, exists := Rooms[joinData.RoomName]
		if !exists {
			log.Println("Room not found:", joinData.RoomName)
			return
		}

		room.AddUser(user)

		msg, err := entities.NewClientMessage(ws, []byte("Joined"))
		if err != nil {
			log.Println(err)
			return
		}
		Messages <- msg
	default:
	}
}
