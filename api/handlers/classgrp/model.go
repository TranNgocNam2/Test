package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"Backend/internal/weekday"
	"github.com/google/uuid"
	"net/mail"
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
		if err != nil || slotStartDate.Before(time.Now().UTC()) {
			return class.NewClass{}, model.ErrInvalidSlotStartDate
		}

		slotStartTime, err = time.Parse(time.TimeOnly, newClassRequest.Slots.StartTime)
		if err != nil {
			return class.NewClass{}, model.ErrInvalidTime
		}
	}

	newClass := class.NewClass{
		ID:        uuid.New(),
		ProgramId: programId,
		SubjectId: subjectId,
		Name:      newClassRequest.Name,
		Code:      newClassRequest.Code,
		Link:      &newClassRequest.Link,
		Type:      int16(*newClassRequest.Type),
		Password:  newClassRequest.Password,
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

func toCoreUpdateClass(updateClassRequest payload.UpdateClass) class.UpdateClass {
	updateClass := class.UpdateClass{
		Name:     updateClassRequest.Name,
		Code:     updateClassRequest.Code,
		Password: &updateClassRequest.Password,
		Type:     int16(*updateClassRequest.Type),
	}

	return updateClass
}

func validateUpdateClassRequest(updateClassRequest payload.UpdateClass) error {
	if err := validate.Check(updateClassRequest); err != nil {
		return err
	}
	return nil
}

func validateUpdateSlotRequest(updateSlotRequest payload.UpdateSlots) error {
	if err := validate.Check(updateSlotRequest); err != nil {
		return err
	}
	return nil
}

func toCoreImportLearners(importLearnersRequest payload.ImportLearners) (class.ImportLearners, error) {
	for _, email := range importLearnersRequest.Emails {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return class.ImportLearners{}, model.ErrInvalidEmail
		}

	}

	return class.ImportLearners{
		Emails: importLearnersRequest.Emails,
	}, nil
}

func validateImportLearnersRequest(request payload.ImportLearners) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSlots(updateSlotRequest payload.UpdateSlots) ([]class.UpdateSlot, error) {
	var updateSlots []class.UpdateSlot
	for _, slot := range updateSlotRequest.Slots {
		startTime, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			return nil, model.ErrInvalidTime
		}

		endTime, err := time.Parse(time.RFC3339, slot.EndTime)
		if err != nil {
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
	startTime, err := time.Parse(time.RFC3339, checkTeacherTime.StartTime)
	if err != nil {
		return class.CheckTeacherTime{}, model.ErrInvalidTime
	}

	endTime, err := time.Parse(time.RFC3339, checkTeacherTime.EndTime)
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

func toCoreUpdateMeetingLink(updateMeetingLink payload.UpdateMeetingLink) class.UpdateMeeting {
	updateMeeting := class.UpdateMeeting{
		Link: updateMeetingLink.Link,
	}
	return updateMeeting
}

func validateUpdateMeetingLinkRequest(updateMeetingLink payload.UpdateMeetingLink) error {
	if err := validate.Check(updateMeetingLink); err != nil {
		return err
	}
	return nil
}

func toCoreAddLearner(req payload.AddLearner) class.AddLearner {
	return class.AddLearner{
		LearnerId: req.LearnerId,
	}
}

func validateAddLearnerRequest(req payload.AddLearner) error {
	if err := validate.Check(req); err != nil {
		return err
	}
	return nil
}

func toCoreRemoveLearner(req payload.RemoveLearner) class.RemoveLearner {
	return class.RemoveLearner{
		LearnerId: req.LearnerId,
	}
}

func validateRemoveLearnerRequest(req payload.RemoveLearner) error {
	if err := validate.Check(req); err != nil {
		return err
	}
	return nil
}
