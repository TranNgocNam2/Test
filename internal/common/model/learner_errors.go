package model

import "github.com/pkg/errors"

var (
	ErrUnauthorizedFeatureAccess   = errors.New("Tài khoản của bạn không được phép sử dụng tính năng này!")
	ErrClassStarted                = errors.New("Không thể tham gia lớp học đã bắt đầu!")
	ErrWrongPassword               = errors.New("Mật khẩu lớp học không đúng!")
	ErrAlreadyJoinedSpecialization = errors.New("Bạn đã tham gia chuyên ngành này rồi!")
	LearnerNotInClass              = errors.New("Học viên không tham gia lớp học này!")
	ErrInvalidAttendanceCode       = errors.New("Mã điểm danh không hợp lệ!")
	ErrLearnerAlreadyVerified      = errors.New("Người dùng đã được xác thực, vui lòng không thay đổi thông tin!")
	ErrVerificationNotFound        = errors.New("Không tìm thấy thông tin xác thực!")
)
