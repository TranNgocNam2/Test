package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/common/model"
	"Backend/internal/password"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"Backend/internal/weekday"
	"github.com/google/uuid"
	"time"
)

func toCoreNewClass(newClassRequest payload.NewClass) (class.NewClass, error) {
	subjectId, err := uuid.Parse(newClassRequest.SubjectId)
	if err != nil {
		return class.NewClass{}, model.ErrInvalidSubjectId
	}

	programId, err := uuid.Parse(newClassRequest.ProgramId)
	if err != nil {
		return class.NewClass{}, model.ErrInvalidProgramId
	}

	var slotStartDate time.Time
	var slotStartTime time.Time
	if newClassRequest.Slots.StartTime != "" && newClassRequest.Slots.StartDate != "" {
		slotStartDate, err = time.Parse(time.DateOnly, newClassRequest.Slots.StartDate)
		if err != nil || slotStartDate.Before(time.Now()) {
			return class.NewClass{}, model.ErrInvalidSlotStartDate
		}

		slotStartTime, err = time.Parse(time.TimeOnly, newClassRequest.Slots.StartTime)
		if err != nil {
			return class.NewClass{}, model.ErrInvalidTime
		}
	}

	pwd, err := password.Hash(newClassRequest.Password)
	if err != nil {
		return class.NewClass{}, model.ErrInvalidPassword
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

func validateNewClassRequest(newClassRequest payload.NewClass) error {
	if err := validate.Check(newClassRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateClass(updateClassRequest payload.UpdateClass) (class.UpdateClass, error) {
	updateClass := class.UpdateClass{
		Name: updateClassRequest.Name,
		Code: updateClassRequest.Code,
	}

	if updateClassRequest.Password != "" {
		pwd, err := password.Hash(updateClassRequest.Password)
		if err != nil {
			return class.UpdateClass{}, model.ErrInvalidPassword
		}
		updateClass.Password = &pwd
	}

	return updateClass, nil
}

func validateUpdateClassRequest(updateClassRequest payload.UpdateClass) error {
	if err := validate.Check(updateClassRequest); err != nil {
		return err
	}
	return nil
}

func validateUpdateClassTeacherRequest(updateClassTeacherRequest payload.UpdateClassTeacher) error {
	if err := validate.Check(updateClassTeacherRequest); err != nil {
		return err
	}
	return nil
}

func validateUpdateSlotRequest(updateSlotRequest payload.UpdateSlot) error {
	if err := validate.Check(updateSlotRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSlot(updateSlotRequest payload.UpdateSlot) ([]class.UpdateSlot, error) {
	var updateSlots []class.UpdateSlot
	for _, slot := range updateSlotRequest.Slots {
		startTime, err := time.Parse(time.DateTime, slot.StartTime)
		if err != nil || startTime.Before(time.Now()) {
			return nil, model.ErrInvalidTime
		}

		endTime, err := time.Parse(time.DateTime, slot.EndTime)
		if err != nil || endTime.Before(startTime) {
			return nil, model.ErrInvalidTime
		}

		slotId, err := uuid.Parse(slot.ID)
		if err != nil {
			return nil, model.ErrInvalidSlotId
		}

		updateSlot := class.UpdateSlot{
			ID:        slotId,
			StartTime: startTime,
			EndTime:   endTime,
			TeacherId: slot.TeacherId,
			Index:     int32(slot.Index),
		}

		updateSlots = append(updateSlots, updateSlot)
	}

	return updateSlots, nil
}

func toCoreCheckTeacherTime(checkTeacherTime payload.CheckTeacherTime) (class.CheckTeacherTime, error) {
	startTime, err := time.Parse(time.DateTime, checkTeacherTime.StartTime)
	if err != nil {
		return class.CheckTeacherTime{}, model.ErrInvalidTime
	}

	endTime, err := time.Parse(time.DateTime, checkTeacherTime.EndTime)
	if err != nil {
		return class.CheckTeacherTime{}, model.ErrInvalidTime
	}

	teacherTime := class.CheckTeacherTime{
		TeacherId: checkTeacherTime.TeacherId,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return teacherTime, nil
}

func validateCheckTeacherTimeRequest(checkTeacherTime payload.CheckTeacherTime) error {
	if err := validate.Check(checkTeacherTime); err != nil {
		return err
	}
	return nil
}
