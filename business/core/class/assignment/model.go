package assignment

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	Id         uuid.UUID       `json:"id"`
	ClassId    uuid.UUID       `json:"classId"`
	Question   json.RawMessage `json:"question"`
	Deadline   time.Time       `json:"deadline"`
	Status     int             `json:"status"`
	Type       int             `json:"type"`
	CanOverdue bool            `json:"canOverdue"`
}

type LearnerAssignment struct {
	LearnerId        string          `json:"learnerId"`
	Grade            float32         `json:"grade"`
	Data             json.RawMessage `json:"data"`
	SubmissionStatus int             `json:"submissionStatus"`
	GradingStatus    int             `json:"gradingStatus"`
	Assignment       Assignment      `json:"assignment"`
}

type LearnerAssignmentQuery struct {
	LearnerId        string          `json:"learnerId"`
	Grade            float32         `json:"grade"`
	Data             json.RawMessage `json:"data"`
	SubmissionStatus int             `json:"submissionStatus"`
	GradingStatus    int             `json:"gradingStatus"`
	AssignmentId     uuid.UUID       `json:"assignmentId"`
}
