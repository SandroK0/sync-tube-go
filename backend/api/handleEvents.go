package api

import (
	"encoding/json"
	"fmt"

	"github.com/SandroK0/sync-tube-go/backend/entities"
	"github.com/gorilla/websocket"
)

func getRoom(roomName string) (*entities.Room, error) {
	room, exists := Rooms[roomName]
	if !exists {
		return nil, fmt.Errorf("room not found: %s", roomName)
	}
	return room, nil
}

func HandleEvents(event ClientEvent, ws *websocket.Conn) {

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

		response := NewJoinEventResponse(user.Token, room.Name)

		msg, err := entities.NewClientMessage(ws, response)
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

		response := NewMessageEventResponse(messageData.Username, messageData.Body)
		msg, err := entities.NewRoomMessage(room.Name, response)
		if err != nil {
			HandleEventError(err, "creating room message")
			return
		}

		Messages <- msg
	default:
		HandleEventError(fmt.Errorf("unknown event type: %s", event.EventType), "handling event")

	}
}
