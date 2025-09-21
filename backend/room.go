package main

import "github.com/gorilla/websocket"

type User struct {
	nickname string
	conn     *websocket.Conn
}

type Room struct {
	name    string
	members []User
}
