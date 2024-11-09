package model

import "github.com/pkg/errors"

var (
	ErrSlotEnded      = errors.New("Giờ học đã kết thúc!")
	ErrSlotNotStarted = errors.New("Giờ học chưa bắt đầu!")
)
