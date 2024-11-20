package payload

type ClassAccess struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LearnerAttendance struct {
	Index          *int   `json:"index" validate:"required"`
	AttendanceCode string `json:"attendanceCode" validate:"required,len=6"`
}

type UpdateLearner struct {
	SchoolId   string   `json:"schoolId" validate:"required"`
	ImageLinks []string `json:"image_links" validate:"required"`
	Type       *int16   `json:"type" validate:"required,gte=0,lte=1"`
}
