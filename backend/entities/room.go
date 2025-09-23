package entities

type Room struct {
	Name  string
	Users []*User
}

func (r *Room) AddUser(user *User) {
	r.Users = append(r.Users, user)
}
