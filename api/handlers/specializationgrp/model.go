package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/slice"
	"Backend/internal/validate"
	"fmt"
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
		ID            *uuid.UUID `json:"id,omitempty"`
		Name          *string    `json:"name,omitempty"`
		Image         *string    `json:"image,omitempty"`
		Code          *string    `json:"code,omitempty"`
		LastUpdated   *time.Time `json:"lastUpdated,omitempty"`
		TotalSessions *int16     `json:"totalSessions,omitempty"`
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
			ID            *uuid.UUID
			Name          *string
			Image         *string
			Code          *string
			LastUpdated   time.Time
			TotalSessions *int16
		}, len(subjectIDs))
		for i, id := range subjectIDs {
			specialization.Subjects[i] = &struct {
				ID            *uuid.UUID
				Name          *string
				Image         *string
				Code          *string
				LastUpdated   time.Time
				TotalSessions *int16
			}{
				ID: &id,
			}
		}
	}

	return specialization, nil
}

func validateNewSpecializationRequest(newSpecializationRequest request.NewSpecialization) error {
	if err := validate.Check(newSpecializationRequest); err != nil {
		return fmt.Errorf(validate.ErrValidation.Error(), err)
	}
	return nil
}
