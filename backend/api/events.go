package api

import (
	"encoding/json"
	"fmt"
)

type ClientEventType string

const (
	JoinRoom    ClientEventType = "join_room"
	CreateRoom  ClientEventType = "create_room"
	SendMessage ClientEventType = "send_message"
)

type ServerEventType string

const (
	RoomJoined      ServerEventType = "room_joined"
	RoomCreated     ServerEventType = "room_created"
	MessageReceived ServerEventType = "message_received"
	Error           ServerEventType = "error"
)

type ClientEvent struct {
	EventType ClientEventType `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

func NewClientEvent(msg []byte) (*ClientEvent, error) {
	var event ClientEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

type CommonEventData struct {
	RoomName string `json:"roomName"`
	Username string `json:"username"`
}

type JoinRoomEventData struct {
	CommonEventData
}

type CreateRoomEventData struct {
	CommonEventData
}

type SendMessageEventData struct {
	CommonEventData
	Body string `json:"body"`
}

type ServerEvent struct {
	EventType ServerEventType `json:"eventType"`
	Data      any             `json:"data"`
}

type ErrorEventData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorEvent(code, message string) ServerEvent {
	data := ErrorEventData{Code: code, Message: message}

	return ServerEvent{
		EventType: Error,
		Data:      data,
	}
}

type RoomJoinedEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

func RoomJoinedEvent(token, roomName string) ServerEvent {

	data := RoomJoinedEventData{Token: token,
		RoomName: roomName}

	return ServerEvent{
		EventType: RoomJoined,
		Data:      data,
	}
}

type MessageReceivedEventData struct {
	Username string `json:"username"`
	Body     string `json:"body"`
}

func MessageReceivedEvent(username, body string) ServerEvent {

	data := MessageReceivedEventData{Username: username,
		Body: body}

	return ServerEvent{
		EventType: MessageReceived,
		Data:      data,
	}
}

type RoomCreatedEventData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

func RoomCreatedEvent(roomName, token string) ServerEvent {
	data := RoomCreatedEventData{Token: token, RoomName: roomName}

	return ServerEvent{
		EventType: RoomCreated,
		Data:      data,
	}
}

func UnmarshalClientEventData(event ClientEvent) (any, error) {
	switch event.EventType {
	case JoinRoom:
		var data JoinRoomEventData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return nil, err
		}
		return data, nil
	case SendMessage:
		var data SendMessageEventData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return nil, err
		}
		return data, nil
	case CreateRoom:
		var data CreateRoomEventData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", event.EventType)
	}
}
