package auth

import "time"

type User struct {
	Id        uint
	GroupId   uint
	Username  string
	Email     string
	CreatedAt time.Time
}
