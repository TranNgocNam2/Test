package user

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"net/mail"
	"time"
)

type User struct {
	ID       string  `json:"id"`
	FullName string  `json:"fullName"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone,omitempty"`
	Photo    *string `json:"photo,omitempty"`
	Role     *int16  `json:"role,omitempty"`
	Status   *int32  `json:"status,omitempty"`
}

type Verification struct {
	ID        uuid.UUID `json:"id"`
	Status    int16     `json:"status"`
	Note      *string   `json:"note"`
	ImageLink []string  `json:"imageLink"`
	Type      int16     `json:"type"`
	School    School    `json:"school"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

type Details struct {
	ID         string  `json:"id"`
	FullName   *string `json:"fullName"`
	Email      string  `json:"email"`
	Phone      *string `json:"phone"`
	Role       int16   `json:"role"`
	Photo      *string `json:"photo"`
	School     *School `json:"school,omitempty"`
	Type       *int16  `json:"type,omitempty"`
	Status     int16   `json:"status,omitempty"`
	IsVerified bool    `json:"isVerified"`
}

type NewUser struct {
	ID       string
	Email    mail.Address
	Role     int16
	FullName string
}

type VerifyLearner struct {
	Status int16
	Note   string
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
		ID:         dbUser.ID,
		FullName:   dbUser.FullName,
		Email:      dbUser.Email,
		Role:       dbUser.AuthRole,
		Phone:      dbUser.Phone,
		Status:     int16(dbUser.Status),
		Photo:      dbUser.ProfilePhoto,
		IsVerified: dbUser.IsVerified,
		Type:       dbUser.Type,
	}
}
