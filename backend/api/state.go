package api

import (
	"github.com/SandroK0/chat-rooms/backend/entities"
	"github.com/gorilla/websocket"
)

type UserRoom struct {
	Username string
	RoomName string
}

func NewUserRoom(username, roomname string) *UserRoom {
	return &UserRoom{Username: username, RoomName: roomname}
}

var (
	Clients      = make(map[*websocket.Conn]bool)
	Messages     = make(chan *Message)
	Rooms        = make(map[string]*entities.Room)
	TokenToRooms = make(map[string]*UserRoom)
)
