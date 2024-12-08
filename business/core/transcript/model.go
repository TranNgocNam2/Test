package transcript

import "github.com/google/uuid"

type LearnerTranscriptQuery struct {
	LearnerId    string    `json:"learnerId"`
	TranscriptId uuid.UUID `json:"transcriptId"`
	Grade        float64   `json:"grade"`
	Status       int32     `json:"status"`
}

type Transcript struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Index  int       `json:"index"`
	Grade  float32   `json:"grade"`
	Status int16     `json:"status"`
}
