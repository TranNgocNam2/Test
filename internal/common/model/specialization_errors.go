package model

import "github.com/pkg/errors"

var (
	ErrSpecCodeAlreadyExist = errors.New("Mã chuyên ngành đã tồn tại!")
	ErrSpecNotFound         = errors.New("Chuyên ngành không tồn tại!")
	ErrSubjectIDsInvalid    = errors.New("Mã môn học không hợp lệ!")
)
