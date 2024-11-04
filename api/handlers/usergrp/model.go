package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/enum/role"
	"net/mail"
)

func toCoreNewUser(newUserRequest payload.NewUser) (user.NewUser, error) {
	authRole := *newUserRequest.Role
	if authRole == role.ADMIN {
		return user.NewUser{}, model.ErrUserCannotBeCreated
	}

	emailAddr, err := mail.ParseAddress(newUserRequest.Email)
	if err != nil {
		return user.NewUser{}, model.ErrInvalidEmail
	}

	newUser := user.NewUser{
		ID:    newUserRequest.ID,
		Email: *emailAddr,
		Role:  int16(authRole),
	}

	return newUser, nil
}
func validateNewUserRequest(newUserRequest payload.NewUser) error {
	if err := validate.Check(newUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateUser(updateUserRequest payload.UpdateUser) (user.UpdateUser, error) {
	authRole := *updateUserRequest.Role
	if authRole == role.LEARNER && updateUserRequest.SchoolID == nil {
		return user.UpdateUser{}, model.ErrNilSchool
	}

	schoolID, err := uuid.Parse(*updateUserRequest.SchoolID)
	if err != nil && updateUserRequest.SchoolID != nil {
		return user.UpdateUser{}, model.ErrInvalidSchoolID
	}

	emailAddr, err := mail.ParseAddress(updateUserRequest.Email)
	if err != nil {
		return user.UpdateUser{}, model.ErrInvalidEmail
	}

	if !user.IsValidPhoneNumber(updateUserRequest.Phone) {
		return user.UpdateUser{}, model.ErrInvalidPhoneNumber
	}

	user := user.UpdateUser{
		FullName: updateUserRequest.FullName,
		Email:    *emailAddr,
		Phone:    updateUserRequest.Phone,
		Gender:   int16(*updateUserRequest.Gender),
		Role:     int16(*updateUserRequest.Role),
		Photo:    updateUserRequest.Photo,
		Image:    updateUserRequest.ImageLink,
	}
	if authRole == role.LEARNER {
		user.SchoolID = &schoolID
	}

	return user, nil
}

func validateUpdateUserRequest(updateUserRequest payload.UpdateUser) error {
	if err := validate.Check(updateUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreVerifyUser(verifyUserRequest payload.VerifyUser) (user.VerifyUser, error) {
	status := int32(verifyUserRequest.Status)
	if status != user.Verified && status != user.Failed {
		return user.VerifyUser{}, model.InvalidUserStatus
	}

	return user.VerifyUser{Status: status}, nil
}

func validateVerifyUserRequest(verifyUserRequest payload.VerifyUser) error {
	if err := validate.Check(verifyUserRequest); err != nil {
		return err
	}
	return nil
}
