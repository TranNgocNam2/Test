package transcript

import "github.com/google/uuid"

type LearnerTranscriptQuery struct {
	LearnerId      string    `json:"learnerId"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	TranscriptId   uuid.UUID `json:"transcriptId"`
	TranscriptName string    `json:"transcriptName"`
	Index          int       `json:"index"`
	Grade          float64   `json:"grade"`
	Status         int32     `json:"status"`
}

type Transcript struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Index  int       `json:"index"`
	Grade  float32   `json:"grade"`
	Status int16     `json:"status"`
}
