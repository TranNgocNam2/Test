package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
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
		ID:       newUserRequest.ID,
		Email:    *emailAddr,
		Role:     int16(authRole),
		FullName: newUserRequest.FullName,
	}

	return newUser, nil
}
func validateNewUserRequest(newUserRequest payload.NewUser) error {
	if err := validate.Check(newUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateUser(updateUserRequest payload.UpdateUser) user.UpdateUser {
	return user.UpdateUser{
		FullName: updateUserRequest.FullName,
		Phone:    updateUserRequest.Phone,
		Photo:    updateUserRequest.Photo,
	}
}

func validateUpdateUserRequest(updateUserRequest payload.UpdateUser) error {
	if err := validate.Check(updateUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreVerifyUser(verifyUserRequest payload.VerifyLearner) (user.VerifyLearner, error) {
	if status.Verification(verifyUserRequest.Status) != status.Verified &&
		status.Verification(verifyUserRequest.Status) != status.Failed {
		return user.VerifyLearner{}, model.InvalidUserStatus
	}

	return user.VerifyLearner{Status: verifyUserRequest.Status}, nil
}

func validateVerifyUserRequest(verifyUserRequest payload.VerifyLearner) error {
	if err := validate.Check(verifyUserRequest); err != nil {
		return err
	}
	return nil
}
