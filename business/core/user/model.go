package user

import (
	"github.com/google/uuid"
	"net/mail"
)

type User struct {
	ID       string
	FullName string
	Email    string
	Phone    string
	Gender   int16
	Role     int16
	Photo    string
	School   struct {
		ID   uuid.UUID
		Name string
	}
}

type NewUser struct {
	ID       string
	FullName string
	Email    mail.Address
	Phone    string
	Gender   int
	Role     int
	Photo    string
	SchoolID *uuid.UUID
}
