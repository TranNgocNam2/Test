package model

import "github.com/pkg/errors"

var (
	ErrCertificateNotFound  = errors.New("Không tìm thấy chứng chỉ!")
	ErrCertificateIdInvalid = errors.New("Mã chứng chỉ không hợp lệ!")
)
