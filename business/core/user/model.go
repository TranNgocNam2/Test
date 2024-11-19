package user

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"net/mail"
)

type Details struct {
	ID             string  `json:"id"`
	FullName       *string `json:"fullName"`
	Email          string  `json:"email"`
	Phone          *string `json:"phone"`
	Photo          *string `json:"photo"`
	School         *School `json:"school,omitempty"`
	Type           *int16  `json:"type,omitempty"`
	Status         int16   `json:"status,omitempty"`
	VerifiedStatus *int16  `json:"verifiedStatus,omitempty"`
}

type NewUser struct {
	ID       string
	Email    mail.Address
	Role     int16
	FullName string
}

type VerifyLearner struct {
	Status int16
}

type UpdateUser struct {
	FullName string
	Phone    string
	Photo    string
}

type School struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

func toCoreUser(dbUser sqlc.User) Details {
	return Details{
		ID:       dbUser.ID,
		FullName: dbUser.FullName,
		Email:    dbUser.Email,
		Phone:    dbUser.Phone,
		Status:   int16(dbUser.Status),
		Photo:    dbUser.ProfilePhoto,
	}
}
