package model

import "github.com/pkg/errors"

var (
	ErrUnauthorizedFeatureAccess   = errors.New("Tài khoản của bạn không được phép sử dụng tính năng này!")
	ErrClassStarted                = errors.New("Không thể tham gia lớp học đã bắt đầu!")
	ErrWrongPassword               = errors.New("Mật khẩu lớp học không đúng!")
	ErrAlreadyJoinedSpecialization = errors.New("Bạn đã tham gia chuyên ngành này rồi!")
)
