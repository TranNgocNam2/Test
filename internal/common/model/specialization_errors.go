package model

import "github.com/pkg/errors"

var (
	ErrSpecCodeAlreadyExist    = errors.New("Mã chuyên ngành đã tồn tại!")
	ErrSpecNotFound            = errors.New("Chuyên ngành không tồn tại!")
	ErrSpecIDInvalid           = errors.New("Mã chuyên ngành không hợp lệ!")
	ErrSpecStatusCannotBeDraft = errors.New("Không thể chuyển trạng thái của chuyên ngành về nháp!")
)
