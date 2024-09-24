package user

import (
	"Backend/business/core/school"
)

type User struct {
	ID       string       `json:"id"`
	FullName string       `json:"fullName"`
	Email    string       `json:"email"`
	Phone    string       `json:"phone"`
	Gender   string       `json:"gender"`
	Role     string       `json:"-"`
	Photo    string       `json:"photo"`
	School   *school.Core `json:"school,omitempty"`
}
