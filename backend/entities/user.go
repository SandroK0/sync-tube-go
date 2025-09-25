package entities

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	Name  string
	Conn  *websocket.Conn
	Token string
}

func NewUser(username string, conn *websocket.Conn) *User {
	token := uuid.New().String()
	return &User{Name: username, Conn: conn, Token: token}
}
