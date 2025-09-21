package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Event struct {
	EventType string `json:"eventType"`
	Data      string `json:"data"`
}

func NewEvent(msg []byte) (*Event, error) {

	var event Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &event, nil
}

func HandleEvents(event Event, ws *websocket.Conn) {

	switch event.EventType {
	case "client":
		msg, err := NewClientMessage(ws, []byte(event.Data))
		if err != nil {
			log.Println(err)
			return
		}
		Messages <- msg
	case "global":
		msg, err := NewBroadcastMessage([]byte(event.Data))
		if err != nil {
			log.Println(err)
			return
		}
		Messages <- msg
	default:
	}
}
