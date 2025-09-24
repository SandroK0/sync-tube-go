package entities

import (
	"errors"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	ClientSpecific  MessageType = "client"
	RoomBroadcast   MessageType = "room"
	GlobalBroadcast MessageType = "global"
)

type Message struct {
	Type     MessageType
	Client   *websocket.Conn // For client-specific messages
	RoomName string          // For room-specific messages
	Content  []byte
}

func NewClientMessage(client *websocket.Conn, content []byte) (*Message, error) {
	if client == nil {
		return nil, errors.New("client connection cannot be nil")
	}
	return &Message{
		Type:    ClientSpecific,
		Client:  client,
		Content: content,
	}, nil
}

func NewRoomMessage(roomName string, content []byte) (*Message, error) {
	if roomName == "" {
		return nil, errors.New("roomID cannot be empty")
	}
	return &Message{
		Type:     RoomBroadcast,
		RoomName: roomName,
		Content:  content,
	}, nil
}

func NewBroadcastMessage(content []byte) (*Message, error) {
	return &Message{
		Type:    GlobalBroadcast,
		Content: content,
	}, nil
}
