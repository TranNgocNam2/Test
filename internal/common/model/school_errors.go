package model

import "github.com/pkg/errors"

var (
	ErrSchoolNotFound    = errors.New("Không tìm thấy trường học!")
	ErrInvalidSchoolID   = errors.New("ID trường học không hợp lệ!")
	ErrInvalidDistrictID = errors.New("ID quận/huyện không hợp lệ!")
	ErrInvalidProvinceID = errors.New("ID tỉnh/thành phố không hợp lệ!")
)
