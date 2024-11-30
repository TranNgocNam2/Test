package payload

type Assignment struct {
	Question   interface{} `json:"question"`
	Deadline   string      `json:"deadline" validate:"required"`
	Type       *int        `json:"type" validate:"gte=0,lte=1,required"`
	Status     *int        `json:"status" validate:"gte=0,lte=1, required"`
	CanOverdue bool        `json:"canOverdue" validate:"required"`
}

type AssignmentGrade struct {
	Grade float32 `json:"grade" validate:"required"`
}
