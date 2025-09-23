package entities

import "github.com/gorilla/websocket"

type User struct {
	Name string
	Conn *websocket.Conn
}

func NewUser(username string, conn *websocket.Conn) *User {
	return &User{Name: username, Conn: conn}
}
