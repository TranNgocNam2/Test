package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
)

func toCoreNewSpecialization(newSpecialization payload.NewSpecialization) specialization.NewSpecialization {
	return specialization.NewSpecialization{
		ID:          uuid.New(),
		Name:        newSpecialization.Name,
		Code:        newSpecialization.Code,
		Description: &newSpecialization.Description,
		TimeAmount:  &newSpecialization.TimeAmount,
		Image:       &newSpecialization.Image,
	}
}

func validateNewSpecializationRequest(newSpecializationRequest payload.NewSpecialization) error {
	if err := validate.Check(newSpecializationRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdatedSpecialization(updateSpecialization payload.UpdateSpecialization) (specialization.UpdateSpecialization, error) {

	var specSubjects []specialization.SpecSubject
	for _, subject := range updateSpecialization.Subjects {
		subjectId, err := uuid.Parse(subject.ID)
		if err != nil {
			return specialization.UpdateSpecialization{}, model.ErrSubjectIDInvalid
		}

		specSubject := specialization.SpecSubject{
			ID:    subjectId,
			Index: int16(*subject.Index),
		}
		specSubjects = append(specSubjects, specSubject)
	}

	specialization := specialization.UpdateSpecialization{
		Name:        updateSpecialization.Name,
		Code:        updateSpecialization.Code,
		Status:      *updateSpecialization.Status,
		Description: updateSpecialization.Description,
		TimeAmount:  updateSpecialization.TimeAmount,
		Image:       updateSpecialization.Image,
		Subjects:    specSubjects,
	}
	return specialization, nil
}
func validateUpdateSpecializationRequest(updateSpecializationRequest payload.UpdateSpecialization) error {
	if err := validate.Check(updateSpecializationRequest); err != nil {
		return err
	}
	return nil
}
