package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/validate"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/mail"
)

var (
	ErrInvalidSchoolID    = errors.New("ID trường học không hợp lệ!")
	ErrInvalidEmail       = errors.New("Email không hợp lệ!")
	ErrInvalidPhoneNumber = errors.New("Số điện thoại không hợp lệ!")
)

type UserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int16  `json:"gender"`
	Role     int16  `json:"-"`
	Photo    string `json:"photo"`
	School   struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"school,omitempty"`
}

func toUserResponse(user user.User) UserResponse {
	userResponse := UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Role:     user.Role,
		Photo:    user.Photo,
	}
	if user.School.ID != uuid.Nil {
		userResponse.School.ID = user.School.ID.String()
		userResponse.School.Name = user.School.Name
	}

	return userResponse
}

func toUserResponses(users []user.User) []UserResponse {
	items := make([]UserResponse, len(users))
	for i, user := range users {
		items[i] = toUserResponse(user)
	}
	return items
}

type NewUserRequest struct {
	ID       string `json:"id" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,startswith=0,len=10"`
	Gender   int    `json:"gender" validate:"required,gte=1,lte=3"`
	Role     int    `json:"role" validate:"required,gte=1,lte=3"`
	Photo    string `json:"photo" validate:"required"`
	SchoolID string `json:"schoolID"`
}

func toCoreNewUser(newUserRequest NewUserRequest) (user.NewUser, error) {
	schoolID, err := uuid.Parse(newUserRequest.SchoolID)
	if err != nil && newUserRequest.SchoolID != "" {
		return user.NewUser{}, ErrInvalidSchoolID
	}

	emailAddr, err := mail.ParseAddress(newUserRequest.Email)
	if err != nil {
		return user.NewUser{}, ErrInvalidEmail
	}

	if !user.IsValidPhoneNumber(newUserRequest.Phone) {
		return user.NewUser{}, ErrInvalidPhoneNumber
	}

	user := user.NewUser{
		ID:       newUserRequest.ID,
		FullName: newUserRequest.FullName,
		Email:    *emailAddr,
		Phone:    newUserRequest.Phone,
		Gender:   newUserRequest.Gender,
		Role:     newUserRequest.Role,
		Photo:    newUserRequest.Photo,
		SchoolID: &schoolID,
	}

	return user, nil
}
func (newUserRequest NewUserRequest) Validate() error {
	if err := validate.Check(newUserRequest); err != nil {
		return fmt.Errorf(validate.ErrValidation.Error(), err)
	}
	return nil
}
