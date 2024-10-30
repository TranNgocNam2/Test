package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/password"
	"Backend/internal/validate"
	"Backend/internal/weekday"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrInvalidSubjectId     = errors.New("Mã môn học không hợp lệ!")
	ErrInvalidProgramId     = errors.New("Mã chương trình học không hợp lệ!")
	ErrInvalidSlotId        = errors.New("Mã buổi học không hợp lệ!")
	ErrInvalidStartDate     = errors.New("Thời gian bắt đầu lớp học không hợp lệ!")
	ErrInvalidSlotStartDate = errors.New("Thời gian bắt đầu buổi học không hợp lệ!")
	ErrInvalidPassword      = errors.New("Mật khẩu không hợp lệ!")
	ErrInvalidTime          = errors.New("Thời gian không hợp lệ!")
)

type NewClass struct {
	ProgramId string `json:"programID" validate:"required"`
	SubjectId string `json:"subjectID" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Code      string `json:"code" validate:"required"`
	Link      string `json:"link"`
	Slots     struct {
		WeekDays  []int  `json:"weekDays" validate:"gte=0,lte=6"`
		StartTime string `json:"startTime"`
		StartDate string `json:"startDate"`
	} `json:"slots"`
	Password string `json:"password" validate:"required"`
}

func toCoreNewClass(newClassRequest NewClass) (class.NewClass, error) {
	subjectId, err := uuid.Parse(newClassRequest.SubjectId)
	if err != nil {
		return class.NewClass{}, ErrInvalidSubjectId
	}

	programId, err := uuid.Parse(newClassRequest.ProgramId)
	if err != nil {
		return class.NewClass{}, ErrInvalidProgramId
	}

	var slotStartDate time.Time
	var slotStartTime time.Time
	if newClassRequest.Slots.StartTime != "" && newClassRequest.Slots.StartDate != "" {
		slotStartDate, err = time.Parse(time.DateOnly, newClassRequest.Slots.StartDate)
		if err != nil || slotStartDate.Before(time.Now()) {
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
		ProgramId: programId,
		SubjectId: subjectId,
		Name:      newClassRequest.Name,
		Link:      &newClassRequest.Link,
		Code:      newClassRequest.Code,
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

type UpdateClass struct {
	Name     string `json:"name" validate:"required"`
	Link     string `json:"link" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func toCoreUpdateClass(updateClassRequest UpdateClass) (class.UpdateClass, error) {
	pwd, err := password.Hash(updateClassRequest.Password)
	if err != nil {
		return class.UpdateClass{}, ErrInvalidPassword
	}

	updateClass := class.UpdateClass{
		Name:     updateClassRequest.Name,
		Link:     updateClassRequest.Link,
		Password: pwd,
	}

	return updateClass, nil
}

func validateUpdateClassRequest(updateClassRequest UpdateClass) error {
	if err := validate.Check(updateClassRequest); err != nil {
		return err
	}
	return nil
}
func validateNewClassRequest(newClassRequest NewClass) error {
	if err := validate.Check(newClassRequest); err != nil {
		return err
	}
	return nil
}

type UpdateClassTeacher struct {
	TeacherIds []string `json:"teacherIds" validate:"required"`
}

func validateUpdateClassTeacherRequest(updateClassTeacherRequest UpdateClassTeacher) error {
	if err := validate.Check(updateClassTeacherRequest); err != nil {
		return err
	}
	return nil
}

type UpdateSlot struct {
	Slots []struct {
		ID        string `json:"id" validate:"required"`
		StartTime string `json:"startTime" validate:"required"`
		EndTime   string `json:"endTime" validate:"required"`
		TeacherId string `json:"teacherId" validate:"required"`
		Index     int    `json:"index" validate:"required"`
	} `json:"slots" validate:"required"`
}

func validateUpdateSlotRequest(updateSlotRequest UpdateSlot) error {
	if err := validate.Check(updateSlotRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSlot(updateSlotRequest UpdateSlot) ([]class.UpdateSlot, error) {
	var updateSlots []class.UpdateSlot
	for _, slot := range updateSlotRequest.Slots {
		startTime, err := time.Parse(time.DateTime, slot.StartTime)
		if err != nil || startTime.Before(time.Now()) {
			return nil, ErrInvalidTime
		}

		endTime, err := time.Parse(time.DateTime, slot.EndTime)
		if err != nil || endTime.Before(startTime) {
			return nil, ErrInvalidTime
		}

		slotId, err := uuid.Parse(slot.ID)
		if err != nil {
			return nil, ErrInvalidSlotId
		}

		updateSlot := class.UpdateSlot{
			ID:        slotId,
			StartTime: startTime,
			EndTime:   endTime,
			TeacherId: slot.TeacherId,
		}

		updateSlots = append(updateSlots, updateSlot)
	}

	return updateSlots, nil
}
