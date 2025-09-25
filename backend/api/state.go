package api

import (
	"github.com/SandroK0/sync-tube-go/backend/entities"
	"github.com/gorilla/websocket"
)

var (
	Clients  = make(map[*websocket.Conn]bool)
	Messages = make(chan *Message)
	Rooms    = make(map[string]*entities.Room)
)
