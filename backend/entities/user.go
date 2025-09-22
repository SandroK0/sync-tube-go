package entities

import "github.com/gorilla/websocket"

type User struct {
	Nickname string
	Conn     *websocket.Conn
}

func NewUser(nickname string, conn *websocket.Conn) *User {
	return &User{Nickname: nickname, Conn: conn}
}
