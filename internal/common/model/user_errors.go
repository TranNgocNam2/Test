package model

import "github.com/pkg/errors"

var (
	ErrInvalidEmail            = errors.New("Email không hợp lệ!")
	ErrInvalidPhoneNumber      = errors.New("Số điện thoại không hợp lệ!")
	ErrUserCannotBeCreated     = errors.New("Không thể tạo người dùng!")
	ErrNilSchool               = errors.New("Vui lòng cung cấp thông tin về trường học!")
	InvalidUserStatus          = errors.New("Trạng thái người dùng không hợp lệ!")
	ErrEmailAlreadyExists      = errors.New("Email đã tồn tại trong hệ thống!")
	ErrPhoneAlreadyExists      = errors.New("Số điện thoại đã tồn tại trong hệ thống!")
	ErrUserAlreadyExist        = errors.New("Người dùng đã tồn tại trong hệ thống!")
	ErrUserNotFound            = errors.New("Người dùng không tồn tại trong hệ thống!")
	ErrUserCannotBeVerified    = errors.New("Người dùng không thể được xác thực!")
	ErrInvalidVerificationInfo = errors.New("Thông tin xác thực của người dùng không hợp lệ!")
)
