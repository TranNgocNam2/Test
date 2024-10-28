package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/password"
	"Backend/internal/validate"
	"Backend/internal/weekday"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"time"
)

var (
	ErrInvalidSubjectID     = errors.New("Mã môn học không hợp lệ!")
	ErrInvalidProgramID     = errors.New("Mã chương trình học không hợp lệ!")
	ErrInvalidStartDate     = errors.New("Thời gian bắt đầu lớp học không hợp lệ!")
	ErrInvalidSlotStartDate = errors.New("Thời gian bắt đầu buổi học không hợp lệ!")
	ErrInvalidPassword      = errors.New("Mật khẩu không hợp lệ!")
	ErrInvalidWeekDay       = errors.New("Ngày trong tuần không hợp lệ!")
	ErrInvalidTime          = errors.New("Thời gian không hợp lệ!")
)

func toCoreNewClass(newClassRequest request.NewClass) (class.NewClass, error) {
	subjectID, err := uuid.Parse(newClassRequest.SubjectID)
	if err != nil {
		return class.NewClass{}, ErrInvalidSubjectID
	}

	programID, err := uuid.Parse(newClassRequest.ProgramID)
	if err != nil {
		return class.NewClass{}, ErrInvalidProgramID
	}

	startDate, err := time.Parse(time.DateOnly, newClassRequest.StartDate)
	if err != nil || startDate.Before(time.Now()) {
		return class.NewClass{}, ErrInvalidStartDate
	}

	var slotStartDate time.Time
	var slotStartTime time.Time
	if newClassRequest.Slots.StartTime != "" && newClassRequest.Slots.StartDate == "" {
		slotStartDate, err = time.Parse(time.DateOnly, newClassRequest.Slots.StartDate)
		if err != nil || slotStartDate.Before(startDate) {
			return class.NewClass{}, ErrInvalidSlotStartDate
		}

		slotStartTime, err = time.Parse(time.TimeOnly, newClassRequest.Slots.StartTime)
		if err != nil {
			return class.NewClass{}, ErrInvalidTime
		}
	}

	pwd, err := password.Hash(newClassRequest.Password)
	if err != nil {
		return class.NewClass{}, ErrInvalidPassword
	}

	newClass := class.NewClass{
		ID:        uuid.New(),
		ProgramID: programID,
		SubjectID: subjectID,
		Name:      newClassRequest.Name,
		Link:      &newClassRequest.Link,
		Code:      newClassRequest.Code,
		StartDate: startDate,
		Password:  pwd,
	}

	newClass.Slots = struct {
		WeekDays  []time.Weekday
		StartTime *time.Time
		StartDate *time.Time
	}{
		WeekDays:  weekday.Parse(newClassRequest.Slots.WeekDays),
		StartTime: &slotStartTime,
		StartDate: &slotStartDate,
	}
	return newClass, nil
}

func validateNewClassRequest(newClassRequest request.NewClass) error {
	if err := validate.Check(newClassRequest); err != nil {
		return err
	}
	return nil
}
