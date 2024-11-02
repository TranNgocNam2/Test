package model

import "errors"

var (
	ErrSubjectIDInvalid        = errors.New("ID môn học không hợp lệ!")
	ErrSkillNotFound           = errors.New("Kỹ năng không có trong hệ thống!")
	ErrSubjectNotFound         = errors.New("Môn học không có trong hệ thống!")
	ErrCodeAlreadyExist        = errors.New("Mã môn đã tồn tại!")
	ErrSkillRequired           = errors.New("Môn học cần có ít nhất một kĩ năng!")
	ErrInvalidSkillId          = errors.New("Skill id không phải định dạng uuid!")
	ErrInvalidSessions         = errors.New("Số lượng session cho môn học không hợp lệ!")
	ErrInvalidMaterials        = errors.New("Buổi học phải có ít nhất 1 nội dung!")
	ErrInvalidMaterialType     = errors.New("Material có type không phù hợp!")
	ErrDataConversion          = errors.New("Không convert qua json được!")
	ErrInvalidTranscript       = errors.New("Phải có ít nhất 1 cột điểm!")
	ErrInvalidTranscriptWeight = errors.New("Tổng các cột điểm phải chiếm đủ 100%!")
)
