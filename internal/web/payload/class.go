package payload

type NewClass struct {
	ProgramId string `json:"programId" validate:"required"`
	SubjectId string `json:"subjectId" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Code      string `json:"code" validate:"required"`
	Link      string `json:"link"`
	Slots     struct {
		WeekDays  []int  `json:"weekDays" validate:"gte=0,lte=6"`
		StartTime string `json:"startTime"`
		StartDate string `json:"startDate"`
	} `json:"slots"`
	Password string `json:"password" validate:"required"`
}

type UpdateClass struct {
	Name     string `json:"name" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Password string `json:"password"`
}

type UpdateClassTeacher struct {
	TeacherIds []string `json:"teacherIds" validate:"required"`
}

type UpdateSlot struct {
	Slots []struct {
		ID        string `json:"id" validate:"required"`
		StartTime string `json:"startTime" validate:"required"`
		EndTime   string `json:"endTime" validate:"required"`
		TeacherId string `json:"teacherId"`
		Index     int    `json:"index" validate:"required"`
	} `json:"slots" validate:"required"`
}

type CheckTeacherTime struct {
	TeacherId string `json:"teacherId" validate:"required"`
	ClassId   string `json:"classId" validate:"required"`
}
