package learnergrp

import (
	"Backend/business/core/learner"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
)

func toCoreUpdateLearner(request payload.UpdateLearner) (learner.UpdateLearner, error) {
	schoolId, err := uuid.Parse(request.SchoolId)
	if err != nil {
		return learner.UpdateLearner{}, model.ErrInvalidSchoolID
	}

	return learner.UpdateLearner{
		SchoolId:   schoolId,
		ImageLinks: request.ImageLinks,
		Type:       *request.Type,
	}, nil
}

func validateUpdateLearnerRequest(request payload.UpdateLearner) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}

func toCoreClassAccess(classAccess payload.ClassAccess) learner.ClassAccess {
	return learner.ClassAccess{
		Code:     classAccess.Code,
		Password: classAccess.Password,
	}
}

func validateNewClassAccessRequest(request payload.ClassAccess) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}

func toCoreSubmitAttendance(request payload.LearnerAttendance) learner.AttendanceSubmission {
	return learner.AttendanceSubmission{
		Index:          int32(*request.Index),
		AttendanceCode: request.AttendanceCode,
	}
}

func validateLearnerAttendanceRequest(request payload.LearnerAttendance) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
