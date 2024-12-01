package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/common/model"
	"Backend/internal/common/status"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"net/mail"
)

func toCoreNewUser(newUserRequest payload.NewUser) (user.NewUser, error) {
	authRole := *newUserRequest.Role
	//if authRole == role.ADMIN {
	//	return user.NewUser{}, model.ErrUserCannotBeCreated
	//}

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
		Photo:    updateUserRequest.Photo,
	}
}

func validateUpdateUserRequest(updateUserRequest payload.UpdateUser) error {
	if err := validate.Check(updateUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreVerifyUser(req payload.VerifyLearner) (user.VerifyLearner, error) {
	if status.Verification(req.Status) == status.Rejected &&
		req.Note == nil {
		return user.VerifyLearner{}, model.InvalidNoteToVerifiedUser
	}

	return user.VerifyLearner{
		Status: req.Status,
		Note:   *req.Note,
	}, nil
}

func validateVerifyUserRequest(verifyUserRequest payload.VerifyLearner) error {
	if err := validate.Check(verifyUserRequest); err != nil {
		return err
	}
	return nil
}

func toCoreCreateLearner(req payload.NewLearner) (user.NewLearner, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return user.NewLearner{}, model.ErrInvalidSchoolID
	}
	emailAddr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return user.NewLearner{}, model.ErrInvalidEmail
	}

	return user.NewLearner{
		ID:       req.ID,
		Email:    emailAddr.String(),
		FullName: req.FullName,
		Type:     int16(*req.Type),
		SchoolID: schoolID,
	}, nil
}

func validateCreateLearnerRequest(req payload.NewLearner) error {
	if err := validate.Check(req); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateLearner(req payload.UpdateLearner) (user.UpdateLearner, error) {
	schoolID, err := uuid.Parse(req.SchoolId)
	if err != nil {
		return user.UpdateLearner{}, model.ErrInvalidSchoolID
	}

	return user.UpdateLearner{
		Type:     int16(*req.Type),
		SchoolID: schoolID,
	}, nil
}

func validateUpdateLearnerRequest(req payload.UpdateLearner) error {
	if err := validate.Check(req); err != nil {
		return err
	}
	return nil
}
