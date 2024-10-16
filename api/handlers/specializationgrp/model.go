package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/slice"
	"Backend/internal/validate"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
	"time"
)

type SpecializationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Status      int16     `json:"status"`
	Description *string   `json:"description"`
	TimeAmount  float64   `json:"timeAmount"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	Skills      []*struct {
		ID   *uuid.UUID `json:"id,omitempty"`
		Name *string    `json:"name,omitempty"`
	} `json:"skills,omitempty"`
	Subjects []*struct {
		ID           *uuid.UUID `json:"id,omitempty"`
		Name         *string    `json:"name,omitempty"`
		Image        *string    `json:"image,omitempty"`
		Code         *string    `json:"code,omitempty"`
		LastUpdated  time.Time  `json:"lastUpdated,omitempty"`
		TotalSession *int64     `json:"totalSession,omitempty"`
	} `json:"subjects,omitempty"`
}

func toCoreNewSpecialization(newSpecialization request.NewSpecialization) (specialization.Specialization, error) {
	skillIDs, err := slice.GetUUIDs(newSpecialization.Skills)
	if err != nil {
		return specialization.Specialization{}, err
	}

	subjectIDs, err := slice.GetUUIDs(newSpecialization.Subjects)
	if err != nil {
		return specialization.Specialization{}, err
	}
	specialization := specialization.Specialization{
		ID:          uuid.New(),
		Name:        newSpecialization.Name,
		Code:        newSpecialization.Code,
		Description: &newSpecialization.Description,
		TimeAmount:  &newSpecialization.TimeAmount,
		Image:       &newSpecialization.Image,
		Skills:      nil,
		Subjects:    nil,
	}
	if skillIDs != nil {
		specialization.Skills = make([]*struct {
			ID   *uuid.UUID
			Name *string
		}, len(skillIDs))

		for i, id := range skillIDs {
			specialization.Skills[i] = &struct {
				ID   *uuid.UUID
				Name *string
			}{
				ID: &id,
			}
		}
	}

	if subjectIDs != nil {
		specialization.Subjects = make([]*struct {
			ID           *uuid.UUID
			Name         *string
			Image        *string
			Code         *string
			LastUpdated  time.Time
			TotalSession *int64
		}, len(subjectIDs))
		for i, id := range subjectIDs {
			specialization.Subjects[i] = &struct {
				ID           *uuid.UUID
				Name         *string
				Image        *string
				Code         *string
				LastUpdated  time.Time
				TotalSession *int64
			}{
				ID: &id,
			}
		}
	}

	return specialization, nil
}

func validateNewSpecializationRequest(newSpecializationRequest request.NewSpecialization) error {
	if err := validate.Check(newSpecializationRequest); err != nil {
		return err
	}
	return nil
}

func toResponseSpecialization(specialization specialization.Specialization) SpecializationResponse {
	response := SpecializationResponse{
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
		response.Skills = make([]*struct {
			ID   *uuid.UUID `json:"id,omitempty"`
			Name *string    `json:"name,omitempty"`
		}, len(specialization.Skills))

		for i, skill := range specialization.Skills {
			response.Skills[i] = &struct {
				ID   *uuid.UUID `json:"id,omitempty"`
				Name *string    `json:"name,omitempty"`
			}{
				ID:   skill.ID,
				Name: skill.Name,
			}
		}
	}

	if specialization.Subjects != nil {
		response.Subjects = make([]*struct {
			ID           *uuid.UUID `json:"id,omitempty"`
			Name         *string    `json:"name,omitempty"`
			Image        *string    `json:"image,omitempty"`
			Code         *string    `json:"code,omitempty"`
			LastUpdated  time.Time  `json:"lastUpdated,omitempty"`
			TotalSession *int64     `json:"totalSession,omitempty"`
		}, len(specialization.Subjects))

		for i, subject := range specialization.Subjects {
			response.Subjects[i] = &struct {
				ID           *uuid.UUID `json:"id,omitempty"`
				Name         *string    `json:"name,omitempty"`
				Image        *string    `json:"image,omitempty"`
				Code         *string    `json:"code,omitempty"`
				LastUpdated  time.Time  `json:"lastUpdated,omitempty"`
				TotalSession *int64     `json:"totalSession,omitempty"`
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

	return response
}
