package user

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"net/mail"
)

type User struct {
	ID       string       `json:"id"`
	FullName *string      `json:"fullName"`
	Email    mail.Address `json:"email"`
	Phone    *string      `json:"phone"`
	Gender   *int16       `json:"gender"`
	Role     int16        `json:"role"`
	Photo    *string      `json:"photo"`
	School   *struct {
		ID   *uuid.UUID `json:"id"`
		Name *string    `json:"name"`
	} `json:"school"`
	ImageLink []string `json:"image_links"`
}

type NewUser struct {
	ID    string
	Email mail.Address
	Role  int16
}

type VerifyUser struct {
	Status int32
}

type UpdateUser struct {
	FullName string
	Role     int16
	Email    mail.Address
	Phone    string
	Gender   int16
	Photo    string
	SchoolID *uuid.UUID
	Image    []string
}

func toCoreUser(dbUser sqlc.User) User {
	emailAddr, _ := mail.ParseAddress(dbUser.Email)
	return User{
		ID:        dbUser.ID,
		FullName:  dbUser.FullName,
		Email:     *emailAddr,
		Phone:     dbUser.Phone,
		Gender:    dbUser.Gender,
		Role:      dbUser.AuthRole,
		Photo:     dbUser.ProfilePhoto,
		ImageLink: dbUser.Image,
	}
}
