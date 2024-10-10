package user

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"net/mail"
)

type User struct {
	ID       string
	FullName *string
	Email    mail.Address
	Phone    *string
	Gender   *int16
	Role     int16
	Photo    *string
	School   *struct {
		ID   *uuid.UUID
		Name *string
	}
}

func toCoreUser(dbUser sqlc.User) User {
	emailAddr, _ := mail.ParseAddress(dbUser.Email)
	return User{
		ID:       dbUser.ID,
		FullName: dbUser.FullName,
		Email:    *emailAddr,
		Phone:    dbUser.Phone,
		Gender:   dbUser.Gender,
		Role:     dbUser.AuthRole,
		Photo:    dbUser.ProfilePhoto,
	}
}
