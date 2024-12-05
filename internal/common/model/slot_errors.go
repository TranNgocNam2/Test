package model

import "github.com/pkg/errors"

var (
	ErrSlotEnded             = errors.New("Buổi học đã kết thúc!")
	ErrSlotAlreadyStarted    = errors.New("Buổi học đã bắt đầu!")
	ErrSlotNotStarted        = errors.New("Buổi học chưa bắt đầu!")
	ErrTeacherIsNotInSlot    = errors.New("Giáo viên không được truy cập vào buổi học này!")
	ErrCannotUpdateSlotTime  = errors.New("Không thể cập nhật thời gian buổi học!")
	ErrSlotOrderTimeConflict = errors.New("Thời gian buổi học hiện tại bị trùng với các buổi học trước!")
	ErrStartTimeNotStarted   = errors.New("Thời gian bắt đầu buổi học chưa bắt đầu!")
	ErrSlotTimeConflict      = "Thời gian buổi học hiện tại bị trùng với buổi học %s!"
)
