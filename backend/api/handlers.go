package api

import (
	"fmt"

	"github.com/SandroK0/chat-rooms/backend/entities"
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

	clientEventData, err := UnmarshalClientEventData(event)
	if err != nil {
		HandleEventError(err, "unmarshaling create_room event")
		return
	}

	switch event.EventType {
	case CreateRoom:

		data, ok := clientEventData.(CreateRoomEventData)
		if !ok {
			HandleEventError(fmt.Errorf("invalid event data type for CreateRoom"), "type assertion")
			return
		}

		if data.RoomName == "" || data.Username == "" {
			HandleEventError(fmt.Errorf("missing roomName or username"), "create_room event")
			return
		}

		if _, exists := Rooms[data.RoomName]; exists {

			errorEvent := ErrorEvent("RoomAlreadyExists", "Room name is taken")

			msg, err := NewClientMessage(ws, errorEvent)
			if err != nil {
				HandleEventError(err, "creating client message")
				return
			}

			Messages <- msg
			return
		}

		room := entities.NewRoom(data.RoomName)

		Rooms[data.RoomName] = room

		user := entities.NewUser(data.Username, ws)
		room.AddUser(user)

		TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

		roomCreatedEvent := RoomCreatedEvent(data.RoomName, user.Token)
		msg1, err := NewClientMessage(ws, roomCreatedEvent)
		if err != nil {
			HandleEventError(err, "creating client message for RoomCreatedEvent")
			return
		}

		roomJoinedEvent := RoomJoinedEvent(user.Token, room.Name)
		msg2, err := NewClientMessage(ws, roomJoinedEvent)
		if err != nil {
			HandleEventError(err, "creating client message for RoomJoinedEvent")
			return
		}

		Messages <- msg1
		Messages <- msg2

	case JoinRoom:

		data, ok := clientEventData.(JoinRoomEventData)
		if !ok {
			HandleEventError(fmt.Errorf("invalid event data type for JoinRoom"), "type assertion")
			return
		}

		if data.RoomName == "" || data.Username == "" {
			HandleEventError(fmt.Errorf("missing roomName or username"), "join event")
			return
		}

		room, err := getRoom(data.RoomName)
		if err != nil {
			HandleEventError(err, "join event")
			return
		}

		user := entities.NewUser(data.Username, ws)

		err = room.AddUser(user)
		if err != nil {
			errorEvent := ErrorEvent("UsernameTaken", "User with that name already exists in that room")

			msg, err := NewClientMessage(ws, errorEvent)
			if err != nil {
				HandleEventError(err, "creating client message")
				return
			}

			Messages <- msg
			return
		}

		TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

		roomJoinedEvent := RoomJoinedEvent(user.Token, room.Name)

		msg, err := NewClientMessage(ws, roomJoinedEvent)
		if err != nil {
			HandleEventError(err, "creating client message")
			return
		}

		Messages <- msg
	case LeaveRoom:

		data, ok := clientEventData.(LeaveRoomEventData)
		if !ok {
			HandleEventError(fmt.Errorf("invalid event data type for LeaveRoom"), "type assertion")
			return
		}

		if data.RoomName == "" || data.Username == "" || data.Token == "" {
			HandleEventError(fmt.Errorf("missing roomName or username or token"), "leave event")
			return
		}

		room, err := getRoom(data.RoomName)
		if err != nil {
			HandleEventError(err, "leave event")
			return
		}

		room.RemoveUser(data.Token)

		roomLeftEvent := RoomLeftEvent(data.Token, data.RoomName)

		msg, err := NewClientMessage(ws, roomLeftEvent)
		if err != nil {
			HandleEventError(err, "creating client message")
			return
		}

		Messages <- msg
	case ReconnectRoom:

		data, ok := clientEventData.(ReconnectRoomEventData)
		if !ok {
			HandleEventError(fmt.Errorf("invalid event data type for ReconnectRoom"), "type assertion")
			return
		}

		userRoom, ok := TokenToRooms[data.Token]
		if !ok {
			fmt.Println("Invalid token:", data.Token)
			return
		}

		room, err := getRoom(userRoom.RoomName)
		if err != nil {
			HandleEventError(err, "reconnect_room event")
			return
		}

		user := room.GetUserByToken(data.Token)
		if user == nil {
			errorEvent := ErrorEvent("InvalidToken", "Token is invalid")

			msg, err := NewClientMessage(ws, errorEvent)
			if err != nil {
				HandleEventError(err, "creating client message")
				return
			}

			Messages <- msg
		}

		user.Conn = ws
		roomJoinedEvent := RoomReconnectedEvent(user.Token, room.Name, user.Name)

		msg, err := NewClientMessage(ws, roomJoinedEvent)
		if err != nil {
			HandleEventError(err, "creating client message")
			return
		}

		Messages <- msg

	case SendMessage:

		data, ok := clientEventData.(SendMessageEventData)
		if !ok {
			HandleEventError(fmt.Errorf("invalid event data type for SendMessage"), "type assertion")
			return
		}

		if data.RoomName == "" || data.Username == "" || data.Body == "" {
			HandleEventError(fmt.Errorf("missing roomName, username, or body"), "message event")
			return
		}

		room, err := getRoom(data.RoomName)
		if err != nil {
			HandleEventError(err, "message event")
			return
		}

		messageReceivedEvent := MessageReceivedEvent(data.Username, data.Body)
		msg, err := NewRoomMessage(room.Name, messageReceivedEvent)
		if err != nil {
			HandleEventError(err, "creating room message")
			return
		}

		Messages <- msg
	default:
		HandleEventError(fmt.Errorf("unknown event type: %s", event.EventType), "handling event")

	}
}
