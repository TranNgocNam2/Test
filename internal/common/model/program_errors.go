package model

import "github.com/pkg/errors"

var (
	ErrProgramIDInvalid    = errors.New("Mã chương trình học không hợp lệ!")
	ErrCannotUpdateProgram = errors.New("Không thể cập nhật chương trình học!")
	ErrCannotDeleteProgram = errors.New("Không thể xóa chương trình học!")
	ErrInvalidStartDate    = errors.New("Thời gian bắt đầu khoá học không hợp lệ!")
	ErrInvalidEndDate      = errors.New("Thời gian kết thúc khoá học không hợp lệ!")
)
