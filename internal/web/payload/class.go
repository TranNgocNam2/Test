package payload

type NewClass struct {
	ProgramId string `json:"programId" validate:"required"`
	SubjectId string `json:"subjectId" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Code      string `json:"code" validate:"required"`
	Link      string `json:"link"`
	Type      *int   `json:"type" validate:"required,gte=0,lte=1"`
	Slots     struct {
		WeekDays  []int  `json:"weekDays" validate:"gte=0,lte=6"`
		StartTime string `json:"startTime"`
		StartDate string `json:"startDate"`
	} `json:"slots"`
	Password string `json:"password" validate:"required,min=6,max=10"`
}

type UpdateClass struct {
	Name     string `json:"name" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Password string `json:"password"`
	Type     *int   `json:"type" validate:"required,gte=0,lte=1"`
}

type UpdateMeetingLink struct {
	Link string `json:"link" validate:"required"`
}

type UpdateSlot struct {
	Status *int `json:"status" validate:"required,gte=0,lte=1"`
	Slots  []struct {
		ID        string `json:"id" validate:"required"`
		StartTime string `json:"startTime" validate:"required"`
		EndTime   string `json:"endTime" validate:"required"`
		TeacherId string `json:"teacherId" validate:"required"`
		Index     int    `json:"index" validate:"required"`
	} `json:"slots" validate:"required"`
}

type CheckTeacherTime struct {
	TeacherId string `json:"teacherId" validate:"required"`
	SlotId    string `json:"slotId" validate:"required"`
	StartTime string `json:"startTime" validate:"required"`
	EndTime   string `json:"endTime" validate:"required"`
}

type ImportLearners struct {
	Emails []string `json:"emails" validate:"required,unique"`
}

type AddLearner struct {
	LearnerId string `json:"learnerId" validate:"required"`
}

type RemoveLearner struct {
	LearnerId string `json:"learnerId" validate:"required"`
}
