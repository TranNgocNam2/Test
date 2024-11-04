package model

import (
	"github.com/pkg/errors"
)

var (
	ErrSubjectIDInvalid        = errors.New("ID môn học không hợp lệ!")
	ErrSubjectNotFound         = errors.New("Môn học không có trong hệ thống!")
	ErrCodeAlreadyExist        = errors.New("Mã môn đã tồn tại!")
	ErrSkillRequired           = errors.New("Môn học cần có ít nhất một kĩ năng!")
	ErrInvalidSessions         = errors.New("Số lượng session cho môn học không hợp lệ!")
	ErrInvalidMaterials        = errors.New("Buổi học phải có ít nhất 1 nội dung!")
	ErrInvalidMaterialType     = errors.New("Material có type không phù hợp!")
	ErrDataConversion          = errors.New("Không convert qua json được!")
	ErrInvalidTranscript       = errors.New("Phải có ít nhất 1 cột điểm!")
	ErrInvalidTranscriptWeight = errors.New("Tổng các cột điểm phải chiếm đủ 100%!")
	ErrSubjectIDsInvalid       = errors.New("Mã môn học không hợp lệ!")
)
