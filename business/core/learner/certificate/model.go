package certificate

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type Certificate struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	CreatedAt      time.Time       `json:"createdAt"`
	Specialization *Specialization `json:"specialization,omitempty"`
	Subject        *Subject        `json:"subject,omitempty"`
	Learner        Learner         `json:"learner,omitempty"`
}

type Learner struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

type Specialization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	TimeAmount  float32   `json:"timeAmount"`
	ImageLink   string    `json:"imageLink"`
	Description string    `json:"description"`
}

type Subject struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	ImageLink   string    `json:"imageLink"`
	Program     *Program  `json:"program,omitempty"`
}

type Program struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func toCoreLearner(learner sqlc.User) Learner {
	return Learner{
		ID:       learner.ID,
		FullName: *learner.FullName,
		Email:    learner.Email,
	}
}

func toCoreProgram(program sqlc.Program) *Program {
	return &Program{
		ID:        program.ID,
		Name:      program.Name,
		StartDate: program.StartDate,
		EndDate:   program.EndDate,
	}
}

func toCoreSubject(subject sqlc.Subject, program sqlc.Program) *Subject {
	return &Subject{
		ID:          subject.ID,
		Name:        subject.Name,
		Code:        subject.Code,
		Description: *subject.Description,
		ImageLink:   *subject.ImageLink,
		Program:     toCoreProgram(program),
	}
}

func toCoreSpecialization(specialization sqlc.Specialization) *Specialization {
	return &Specialization{
		ID:          specialization.ID,
		Name:        specialization.Name,
		Code:        specialization.Code,
		TimeAmount:  *specialization.TimeAmount,
		ImageLink:   *specialization.ImageLink,
		Description: *specialization.Description,
	}
}

func toCoreCertificate(certificate sqlc.Certificate) Certificate {
	return Certificate{
		ID:        certificate.ID,
		Name:      certificate.Name,
		CreatedAt: certificate.CreatedAt,
	}
}
