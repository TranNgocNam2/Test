package payload

type ClassAccess struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LearnerAttendance struct {
	Index          int    `json:"index" validate:"required"`
	AttendanceCode string `json:"attendanceCode" validate:"required,len=6"`
}
