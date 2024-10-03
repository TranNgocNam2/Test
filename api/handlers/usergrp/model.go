package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/validate"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
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
	School   *struct {
		ID   *uuid.UUID `json:"id,omitempty"`
		Name *string    `json:"name,omitempty"`
	} `json:"school,omitempty"`
}

func toUserResponse(user user.User) UserResponse {
	userResponse := UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email.Address,
		Phone:    user.Phone,
		Gender:   user.Gender,
		Role:     user.Role,
		Photo:    user.Photo,
	}
	if user.School != nil {
		userResponse.School.ID = user.School.ID
		userResponse.School.Name = user.School.Name
	}

	return userResponse
}

func toUsersResponse(users []user.User) []UserResponse {
	items := make([]UserResponse, len(users))
	for i, user := range users {
		items[i] = toUserResponse(user)
	}
	return items
}

func toCoreUser(newUserRequest request.NewUser) (user.User, error) {
	schoolID, err := uuid.Parse(newUserRequest.SchoolID)
	if err != nil && newUserRequest.SchoolID != "" {
		return user.User{}, ErrInvalidSchoolID
	}

	emailAddr, err := mail.ParseAddress(newUserRequest.Email)
	if err != nil {
		return user.User{}, ErrInvalidEmail
	}

	if !user.IsValidPhoneNumber(newUserRequest.Phone) {
		return user.User{}, ErrInvalidPhoneNumber
	}

	user := user.User{
		ID:       newUserRequest.ID,
		FullName: newUserRequest.FullName,
		Email:    *emailAddr,
		Phone:    newUserRequest.Phone,
		Gender:   int16(newUserRequest.Gender),
		Role:     int16(newUserRequest.Role),
		Photo:    newUserRequest.Photo,
	}
	if schoolID != uuid.Nil {
		user.School.ID = &schoolID
	}

	return user, nil
}
func validateNewUserRequest(newUserRequest request.NewUser) error {
	if err := validate.Check(newUserRequest); err != nil {
		return fmt.Errorf(validate.ErrValidation.Error(), err)
	}
	return nil
}
