package api

import (
	"encoding/json"
	"fmt"

	"github.com/SandroK0/sync-tube-go/backend/entities"
	"github.com/gorilla/websocket"
)

type EventType string

const (
	JoinEvent    EventType = "join"
	MessageEvent EventType = "message"
)

type Event struct {
	EventType EventType       `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

type CommonEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

type JoinEventData struct {
	CommonEventData
}

type MessageEventData struct {
	CommonEventData
	Body string `json:"body"`
}

func getRoom(roomName string) (*entities.Room, error) {
	room, exists := Rooms[roomName]
	if !exists {
		return nil, fmt.Errorf("room not found: %s", roomName)
	}
	return room, nil
}

func NewEvent(msg []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func HandleEvents(event Event, ws *websocket.Conn) {

	jsonData, err := json.Marshal(event)
	if err != nil {
		HandleEventError(err, "marshaling event")
		return
	}

	switch event.EventType {
	case JoinEvent:

		var joinData JoinEventData
		if err := json.Unmarshal(event.Data, &joinData); err != nil {
			HandleEventError(err, "unmarshaling join event")
			return
		}

		if joinData.RoomName == "" || joinData.Username == "" {
			HandleEventError(fmt.Errorf("missing roomName or username"), "join event")
			return
		}

		room, err := getRoom(joinData.RoomName)
		if err != nil {
			HandleEventError(err, "join event")
			return
		}

		user := entities.NewUser(joinData.Username, ws)
		room.AddUser(user)

		msg, err := entities.NewClientMessage(ws, jsonData)
		if err != nil {
			HandleEventError(err, "creating client message")
			return
		}

		Messages <- msg

	case MessageEvent:
		var messageData MessageEventData
		if err := json.Unmarshal(event.Data, &messageData); err != nil {
			HandleEventError(err, "unmarshaling message event")
			return
		}

		if messageData.RoomName == "" || messageData.Username == "" || messageData.Body == "" {
			HandleEventError(fmt.Errorf("missing roomName, username, or body"), "message event")
			return
		}

		room, err := getRoom(messageData.RoomName)
		if err != nil {
			HandleEventError(err, "message event")
			return
		}

		msg, err := entities.NewRoomMessage(room.Name, jsonData)
		if err != nil {
			HandleEventError(err, "creating room message")
			return
		}

		Messages <- msg
	default:
		HandleEventError(fmt.Errorf("unknown event type: %s", event.EventType), "handling event")

	}
}
