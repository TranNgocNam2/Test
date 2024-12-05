package payload

import "github.com/google/uuid"

type UpdateLearnerTranscript struct {
	Learners []LearnerTranscript `json:"learners" validate:"gt=0, dive, required"`
}

type LearnerTranscript struct {
	LearnerId    string    `json:"learnerId"`
	TranscriptId uuid.UUID `json:"transcriptId" validate:"required"`
	Grade        float32   `json:"grade" validate:"required"`
}
