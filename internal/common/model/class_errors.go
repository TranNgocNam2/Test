package model

import "github.com/pkg/errors"

var (
	ErrClassIdInvalid        = errors.New("Mã lớp học không hợp lệ!")
	ErrProgramNotFound       = errors.New("Không tìm thấy chương trình học!")
	ErrInvalidClassStartTime = errors.New("Thời gian bắt đầu lớp học không hợp lệ!")
	ErrInvalidSlotStartTime  = errors.New("Thời gian bắt đầu buổi học không hợp lệ!")
	ErrInvalidSlotEndTime    = errors.New("Thời gian kết thúc buổi học không hợp lệ!")
	ErrSessionNotFound       = errors.New("Không có buổi học nào trong môn học này!")
	ErrInvalidWeekDay        = errors.New("Số ngày học trong tuần không khớp với số buổi học trong môn học!")
	ErrClassNotFound         = errors.New("Không tìm thấy lớp học!")
	ErrClassCodeAlreadyExist = errors.New("Mã của lớp học đã tồn tại!")
	ErrTeacherNotFound       = errors.New("Không tìm thấy giáo viên!")
	ErrInvalidSlotCount      = errors.New("Số lượng buổi học không hợp lệ!")
	ErrInvalidSlotTime       = errors.New("Thời gian buổi học không hợp lệ!")
	ErrSlotNotFound          = errors.New("Không tìm thấy buổi học!")
	ErrTeacherNotAvailable   = errors.New("Giáo viên không thể dạy vào thời gian này!")
	ErrTeacherIsNotInClass   = errors.New("Giáo viên không thuộc lớp học này!")
	ErrInvalidSubjectId      = errors.New("Mã môn học không hợp lệ!")
	ErrInvalidProgramId      = errors.New("Mã chương trình học không hợp lệ!")
	ErrInvalidSlotId         = errors.New("Mã buổi học không hợp lệ!")
	ErrInvalidSlotStartDate  = errors.New("Thời gian bắt đầu buổi học không hợp lệ!")
	ErrInvalidPassword       = errors.New("Mật khẩu không hợp lệ!")
	ErrInvalidTime           = errors.New("Thời gian không hợp lệ!")
)
