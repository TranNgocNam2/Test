package model

import "github.com/pkg/errors"

var (
	ErrSchoolNotFound    = errors.New("Không tìm thấy trường học!")
	ErrInvalidSchoolID   = errors.New("Mã trường học không hợp lệ!")
	ErrInvalidDistrictID = errors.New("Mã quận/huyện không hợp lệ!")
	ErrInvalidProvinceID = errors.New("Mã tỉnh/thành phố không hợp lệ!")
)
