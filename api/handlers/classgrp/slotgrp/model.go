package slotgrp

import (
	"Backend/business/core/class/slot"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"fmt"
	"time"
)

func toCoreUpdateSlot(req payload.UpdateSlot) (slot.UpdateSlot, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	fmt.Println(startTime)
	if err != nil || startTime.Before(time.Now()) {
		return slot.UpdateSlot{}, model.ErrInvalidTime
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	fmt.Println(endTime)
	if err != nil || endTime.Before(startTime) {
		return slot.UpdateSlot{}, model.ErrInvalidTime
	}
	return slot.UpdateSlot{
		StartTime: startTime,
		EndTime:   endTime,
		TeacherId: req.TeacherId,
	}, nil
}

func validateUpdateSlot(req payload.UpdateSlot) error {
	if err := validate.Check(req); err != nil {
		return err
	}
	return nil
}
