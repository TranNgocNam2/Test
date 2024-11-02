package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrInvalidStartDate = errors.New("Thời gian bắt đầu khoá học không hợp lệ!")
	ErrInvalidEndDate   = errors.New("Thời gian kết thúc khoá học không hợp lệ!")
)

func toCoreNewProgram(newProgramRequest payload.NewProgram) (program.NewProgram, error) {
	startDate, err := time.Parse(time.DateOnly, newProgramRequest.StartDate)
	if err != nil || startDate.Before(time.Now()) {
		return program.NewProgram{}, ErrInvalidStartDate
	}

	endDate, err := time.Parse(time.DateOnly, newProgramRequest.EndDate)
	if err != nil || endDate.Before(startDate) {
		return program.NewProgram{}, ErrInvalidEndDate
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
	if err != nil || startDate.Before(time.Now()) {
		return program.UpdateProgram{}, ErrInvalidStartDate
	}

	endDate, err := time.Parse(time.DateOnly, updateProgramRequest.EndDate)
	if err != nil || endDate.Before(startDate) {
		return program.UpdateProgram{}, ErrInvalidEndDate
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
