package model

import "github.com/pkg/errors"

var (
	ErrInvalidSkillId = errors.New("Mã kỹ năng không hợp lệ!")
	ErrSkillNotFound  = errors.New("Kỹ năng không có trong hệ thống!")
)
