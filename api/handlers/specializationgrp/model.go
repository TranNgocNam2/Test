package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/slice"
	"Backend/internal/validate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"time"
)

var (
	ErrSubjectIDsInvalid = errors.New("ID môn học không hợp lệ!")
	ErrSkillIDsInvalid   = errors.New("ID kỹ năng không hợp lệ!")
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
	Skills      []*struct {
		ID   uuid.UUID `json:"id,omitempty"`
		Name string    `json:"name,omitempty"`
	} `json:"skills,omitempty"`
	Subjects []*struct {
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

	if specialization.Skills != nil {
		specDetailsResponse.Skills = make([]*struct {
			ID   uuid.UUID `json:"id,omitempty"`
			Name string    `json:"name,omitempty"`
		}, len(specialization.Skills))

		for i, skill := range specialization.Skills {
			specDetailsResponse.Skills[i] = &struct {
				ID   uuid.UUID `json:"id,omitempty"`
				Name string    `json:"name,omitempty"`
			}{
				ID:   skill.ID,
				Name: skill.Name,
			}
		}
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
	Skills        []*struct {
		ID   uuid.UUID `json:"id,omitempty"`
		Name string    `json:"name,omitempty"`
	} `json:"skills,omitempty"`
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

	if specialization.Skills != nil {
		specResponse.Skills = make([]*struct {
			ID   uuid.UUID `json:"id,omitempty"`
			Name string    `json:"name,omitempty"`
		}, len(specialization.Skills))

		for i, skill := range specialization.Skills {
			specResponse.Skills[i] = &struct {
				ID   uuid.UUID `json:"id,omitempty"`
				Name string    `json:"name,omitempty"`
			}{
				ID:   skill.ID,
				Name: skill.Name,
			}
		}
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

func toCoreUpdatedSpecialization(updateSpecialization request.UpdateSpecialization) (specialization.UpdateSpecialization, error) {
	skillIDs, err := slice.GetUUIDs(updateSpecialization.Skills)
	if err != nil {
		return specialization.UpdateSpecialization{}, ErrSkillIDsInvalid
	}

	subjectIDs, err := slice.GetUUIDs(updateSpecialization.Subjects)
	if err != nil {
		return specialization.UpdateSpecialization{}, ErrSubjectIDsInvalid
	}

	specialization := specialization.UpdateSpecialization{
		Name:        updateSpecialization.Name,
		Code:        updateSpecialization.Code,
		Status:      *updateSpecialization.Status,
		Description: updateSpecialization.Description,
		TimeAmount:  updateSpecialization.TimeAmount,
		Image:       updateSpecialization.Image,
		Skills:      skillIDs,
		Subjects:    subjectIDs,
	}

	return specialization, nil
}
func validateUpdateSpecializationRequest(updateSpecializationRequest request.UpdateSpecialization) error {
	if err := validate.Check(updateSpecializationRequest); err != nil {
		return err
	}
	return nil
}
