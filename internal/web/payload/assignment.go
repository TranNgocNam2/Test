package payload

type Assignment struct {
	Question   interface{} `json:"question"`
	Deadline   string      `json:"deadline" validate:"required"`
	Type       *int        `json:"type" validate:"required,gte=0,lte=1"`
	Status     *int        `json:"status" validate:"required,gte=0,lte=1"`
	CanOverdue *bool       `json:"canOverdue" validate:"required"`
}

type AssignmentGrade struct {
	Grade float32 `json:"grade" validate:"required"`
}

type LearnerSubmission struct {
	Data interface{} `json:"data"`
}
