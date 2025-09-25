package api

import (
	"encoding/json"
	"fmt"
)

type EventType string

const (
	JoinEvent    EventType = "join"
	MessageEvent EventType = "message"
)

type ClientEvent struct {
	EventType EventType       `json:"eventType"`
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

type JoinEventData struct {
	CommonEventData
}

type MessageEventData struct {
	CommonEventData
	Body string `json:"body"`
}

type ServerEvent struct {
	EventType EventType `json:"eventType"`
	Data      any       `json:"data"`
}

type JoinEventResponseData struct {
	Token    string `json:"token"`
	RoomName string `json:"roomName"`
}

func NewJoinEventResponse(token, roomName string) ServerEvent {

	joinEventResponseData := JoinEventResponseData{Token: token,
		RoomName: roomName}

	return ServerEvent{
		EventType: JoinEvent,
		Data:      joinEventResponseData,
	}
}

type MessageEventResponseData struct {
	Username string `json:"username"`
	Body     string `json:"body"`
}

func NewMessageEventResponse(username, body string) ServerEvent {

	messageEventResponseData := MessageEventResponseData{Username: username,
		Body: body}

	return ServerEvent{
		EventType: MessageEvent,
		Data:      messageEventResponseData,
	}
}

func UnmarshalClientEventData(event ClientEvent) (any, error) {
	switch event.EventType {
	case JoinEvent:
		var data JoinEventData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return nil, err
		}
		return data, nil
	case MessageEvent:
		var data MessageEventData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", event.EventType)
	}
}
