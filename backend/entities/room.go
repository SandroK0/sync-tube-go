package entities

type Room struct {
	Name  string
	Users []*User
}

func NewRoom(roomName string) *Room {
	return &Room{Name: roomName}
}

func (r *Room) AddUser(user *User) {
	r.Users = append(r.Users, user)
}
