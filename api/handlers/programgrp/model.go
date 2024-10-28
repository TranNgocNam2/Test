package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/validate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"time"
)

var (
	ErrInvalidStartDate = errors.New("Thời gian bắt đầu khoá học không hợp lệ!")
	ErrInvalidEndDate   = errors.New("Thời gian kết thúc khoá học không hợp lệ!")
)

type ProgramResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	StartDate     string    `json:"startDate"`
	EndDate       string    `json:"endDate"`
	TotalSubjects int64     `json:"totalSubjects"`
}

func toProgramResponse(program program.Program) ProgramResponse {
	return ProgramResponse{
		ID:            program.ID,
		Name:          program.Name,
		StartDate:     program.StartDate.Format(time.DateOnly),
		EndDate:       program.EndDate.Format(time.DateOnly),
		TotalSubjects: program.TotalClasses,
	}
}

func toCoreProgramsResponse(programs []program.Program) []ProgramResponse {
	programsResponse := make([]ProgramResponse, len(programs))
	for i, program := range programs {
		programsResponse[i] = toProgramResponse(program)
	}
	return programsResponse
}

func toCoreNewProgram(newProgramRequest request.NewProgram) (program.NewProgram, error) {
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

func validateNewProgramRequest(newProgramRequest request.NewProgram) error {
	if err := validate.Check(newProgramRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateProgram(updateProgramRequest UpdateProgram) (program.UpdateProgram, error) {
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

func validateUpdateProgramRequest(updateProgramRequest UpdateProgram) error {
	if err := validate.Check(updateProgramRequest); err != nil {
		return err
	}
	return nil
}

type UpdateProgram struct {
	Name        string `json:"name" validate:"required"`
	StartDate   string `json:"startDate" validate:"required"`
	EndDate     string `json:"endDate" validate:"required"`
	Description string `json:"description" validate:"required"`
}
