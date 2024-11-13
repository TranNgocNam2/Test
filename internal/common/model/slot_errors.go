package model

import "github.com/pkg/errors"

var (
	ErrSlotEnded          = errors.New("Buổi học đã kết thúc!")
	ErrSlotNotStarted     = errors.New("Buổi học chưa bắt đầu!")
	ErrTeacherIsNotInSlot = errors.New("Giáo viên không được truy cập vào buổi học này!")
)
