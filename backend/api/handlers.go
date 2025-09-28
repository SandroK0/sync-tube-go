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

func SendError(Code, Message string, ws *websocket.Conn) {
	errorData := ErrorEventData{Code: Code, Message: Message}
	errorEvent := NewServerEvent(Error, errorData)

	msg, err := NewClientMessage(ws, errorEvent)
	if err != nil {
		HandleEventError(err, "creating client message")
		return
	}

	Messages <- msg
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
			SendError("RoomAlreadyExists", "Room name is taken", ws)
			return
		}

		room := entities.NewRoom(data.RoomName)

		Rooms[data.RoomName] = room

		user := entities.NewUser(data.Username, ws)
		room.AddUser(user)

		TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

		roomCreatedData := RoomCreatedEventData{Token: user.Token, RoomName: data.RoomName}
		roomCreatedEvent := NewServerEvent(RoomCreated, roomCreatedData)
		msg1, err := NewClientMessage(ws, roomCreatedEvent)
		if err != nil {
			HandleEventError(err, "creating client message for RoomCreatedEvent")
			return
		}

		roomJoinedData := RoomJoinedEventData{Token: user.Token, RoomName: room.Name}
		roomJoinedEvent := NewServerEvent(RoomJoined, roomJoinedData)
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
			SendError("UsernameTaken", "User with that name already exists in that room", ws)
			return
		}

		TokenToRooms[user.Token] = NewUserRoom(user.Name, room.Name)

		roomJoinedData := RoomJoinedEventData{Token: user.Token, RoomName: room.Name}
		roomJoinedEvent := NewServerEvent(RoomJoined, roomJoinedData)

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

		roomLeftData := RoomLeftEventData{Token: data.Token, RoomName: data.RoomName}
		roomLeftEvent := NewServerEvent(RoomLeft, roomLeftData)

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
			SendError("InvalidToken", "Token is invalid", ws)
		}

		user.Conn = ws
		roomReconnectedData := RoomReconnectedEventData{Token: user.Token, RoomName: room.Name, Username: user.Name}
		roomJoinedEvent := NewServerEvent(RoomReconnected, roomReconnectedData)

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

		messageReceivedData := MessageReceivedEventData{Username: data.Username, Body: data.Body}
		messageReceivedEvent := NewServerEvent(MessageReceived, messageReceivedData)
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
