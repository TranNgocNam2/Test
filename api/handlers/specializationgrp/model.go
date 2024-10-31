package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/common/model"
	"Backend/internal/slice"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
	"time"
)

type SpecializationDetailsResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Status      int16     `json:"status"`
	Description *string   `json:"description"`
	TimeAmount  float64   `json:"timeAmount"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	Subjects    []*struct {
		ID           uuid.UUID `json:"id,omitempty"`
		Name         string    `json:"name,omitempty"`
		Image        string    `json:"image,omitempty"`
		Code         string    `json:"code,omitempty"`
		LastUpdated  time.Time `json:"lastUpdated,omitempty"`
		TotalSession int64     `json:"totalSession,omitempty"`
	} `json:"subjects,omitempty"`
}

func toResponseSpecializationDetails(specialization specialization.Details) SpecializationDetailsResponse {
	specDetailsResponse := SpecializationDetailsResponse{
		ID:          specialization.ID,
		Name:        specialization.Name,
		Code:        specialization.Code,
		Status:      specialization.Status,
		Description: specialization.Description,
		TimeAmount:  *specialization.TimeAmount,
		Image:       *specialization.Image,
		CreatedAt:   specialization.CreatedAt,
	}

	if specialization.Subjects != nil {
		specDetailsResponse.Subjects = make([]*struct {
			ID           uuid.UUID `json:"id,omitempty"`
			Name         string    `json:"name,omitempty"`
			Image        string    `json:"image,omitempty"`
			Code         string    `json:"code,omitempty"`
			LastUpdated  time.Time `json:"lastUpdated,omitempty"`
			TotalSession int64     `json:"totalSession,omitempty"`
		}, len(specialization.Subjects))

		for i, subject := range specialization.Subjects {
			specDetailsResponse.Subjects[i] = &struct {
				ID           uuid.UUID `json:"id,omitempty"`
				Name         string    `json:"name,omitempty"`
				Image        string    `json:"image,omitempty"`
				Code         string    `json:"code,omitempty"`
				LastUpdated  time.Time `json:"lastUpdated,omitempty"`
				TotalSession int64     `json:"totalSession,omitempty"`
			}{
				ID:           subject.ID,
				Name:         subject.Name,
				Image:        subject.Image,
				Code:         subject.Code,
				LastUpdated:  subject.LastUpdated,
				TotalSession: subject.TotalSession,
			}
		}
	}

	return specDetailsResponse
}

type SpecializationResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Status        int16     `json:"status"`
	Image         string    `json:"image"`
	TotalSubjects int64     `json:"totalSubjects"`
}

func toSpecializationResponse(specialization specialization.Specialization) SpecializationResponse {
	specResponse := SpecializationResponse{
		ID:            specialization.ID,
		Name:          specialization.Name,
		Code:          specialization.Code,
		Status:        specialization.Status,
		Image:         *specialization.Image,
		TotalSubjects: specialization.TotalSubject,
	}

	return specResponse
}

func toSpecializationsResponse(specializations []specialization.Specialization) []SpecializationResponse {
	items := make([]SpecializationResponse, len(specializations))
	for i, specialization := range specializations {
		items[i] = toSpecializationResponse(specialization)
	}
	return items

}

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
