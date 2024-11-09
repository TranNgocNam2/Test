package learnergrp

import (
	"Backend/business/core/learner"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
)

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
