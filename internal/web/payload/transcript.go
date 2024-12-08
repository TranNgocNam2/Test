package payload

import "github.com/google/uuid"

type LearnerTranscript struct {
	LearnerId    string    `json:"learnerId"`
	TranscriptId uuid.UUID `json:"transcriptId" validate:"required"`
	Grade        float32   `json:"grade" validate:"required"`
}
