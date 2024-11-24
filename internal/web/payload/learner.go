package payload

type ClassAccess struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LearnerAttendance struct {
	Index          *int   `json:"index" validate:"required"`
	AttendanceCode string `json:"attendanceCode" validate:"required,len=6"`
}

type UpdateVerificationInfo struct {
	SchoolId   string   `json:"schoolId" validate:"required"`
	ImageLinks []string `json:"imageLinks" validate:"required"`
	Type       *int16   `json:"type" validate:"required,gte=0,lte=1"`
}
