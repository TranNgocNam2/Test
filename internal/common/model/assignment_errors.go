package model

import "errors"

var (
	ErrTimeFormat                   = errors.New("Deadline không đúng định dạng")
	ErrInvalidDeadlineTime          = errors.New("Thời gian deadline không hợp lý")
	InvalidClassAssignment          = errors.New("Assignment không nằm trong class này")
	ErrAssignmentDeletion           = errors.New("Assignment không thể xóa")
	ErrAssignmentNotFound           = errors.New("Không tìm thấy Assignment")
	ErrCannotGradeVisibleAssignment = errors.New("Assignment chưa open")
	ErrLearnerAssignmentNotFound    = errors.New("Không tìm thấy assignment của learner")
	ErrGradingNotStartedAssignment  = errors.New("Assignment này chưa bắt đầu")
	ErrInvalidAssignmentSubmision   = errors.New("Assignment này không gửi bài được")
	ErrInvalidAssignmentId          = errors.New("Assignment Id khong phai uuid")
	ErrSubmitOverdue                = errors.New("Assignment khong chap nhan nop tre")
)
