package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/common/model"
	"Backend/internal/slice"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
)

func toCoreNewSpecialization(newSpecialization request.NewSpecialization) specialization.NewSpecialization {
	return specialization.NewSpecialization{
		ID:          uuid.New(),
		Name:        newSpecialization.Name,
		Code:        newSpecialization.Code,
		Description: &newSpecialization.Description,
		TimeAmount:  &newSpecialization.TimeAmount,
		Image:       &newSpecialization.Image,
	}
}

func validateNewSpecializationRequest(newSpecializationRequest request.NewSpecialization) error {
	if err := validate.Check(newSpecializationRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdatedSpecialization(updateSpecialization payload.UpdateSpecialization) (specialization.UpdateSpecialization, error) {
	subjectIDs, err := slice.GetUUIDs(updateSpecialization.Subjects)
	if err != nil {
		return specialization.UpdateSpecialization{}, model.ErrSubjectIDsInvalid
	}

	specialization := specialization.UpdateSpecialization{
		Name:        updateSpecialization.Name,
		Code:        updateSpecialization.Code,
		Status:      *updateSpecialization.Status,
		Description: updateSpecialization.Description,
		TimeAmount:  updateSpecialization.TimeAmount,
		Image:       updateSpecialization.Image,
		Subjects:    subjectIDs,
	}

	return specialization, nil
}
func validateUpdateSpecializationRequest(updateSpecializationRequest payload.UpdateSpecialization) error {
	if err := validate.Check(updateSpecializationRequest); err != nil {
		return err
	}
	return nil
}
