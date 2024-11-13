package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"time"
)

func toCoreNewProgram(newProgramRequest payload.NewProgram) (program.NewProgram, error) {
	startDate, err := time.Parse(time.DateOnly, newProgramRequest.StartDate)
	if err != nil || startDate.Before(time.Now().UTC()) {
		return program.NewProgram{}, model.ErrInvalidStartDate
	}

	endDate, err := time.Parse(time.DateOnly, newProgramRequest.EndDate)
	if err != nil || endDate.Before(startDate) {
		return program.NewProgram{}, model.ErrInvalidEndDate
	}

	newProgram := program.NewProgram{
		ID:          uuid.New(),
		Name:        newProgramRequest.Name,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: newProgramRequest.Description,
	}

	return newProgram, nil
}

func validateNewProgramRequest(newProgramRequest payload.NewProgram) error {
	if err := validate.Check(newProgramRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateProgram(updateProgramRequest payload.UpdateProgram) (program.UpdateProgram, error) {
	startDate, err := time.Parse(time.DateOnly, updateProgramRequest.StartDate)
	if err != nil || startDate.Before(time.Now().UTC()) {
		return program.UpdateProgram{}, model.ErrInvalidStartDate
	}

	endDate, err := time.Parse(time.DateOnly, updateProgramRequest.EndDate)
	if err != nil || endDate.Before(startDate) {
		return program.UpdateProgram{}, model.ErrInvalidEndDate
	}

	updateProgram := program.UpdateProgram{
		Name:        updateProgramRequest.Name,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: updateProgramRequest.Description,
	}

	return updateProgram, nil
}

func validateUpdateProgramRequest(updateProgramRequest payload.UpdateProgram) error {
	if err := validate.Check(updateProgramRequest); err != nil {
		return err
	}
	return nil
}
