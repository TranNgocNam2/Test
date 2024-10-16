// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Assignment struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	TranscriptID   uuid.UUID  `db:"transcript_id" json:"transcriptId"`
	ClassTeacherID uuid.UUID  `db:"class_teacher_id" json:"classTeacherId"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy      *string    `db:"updated_by" json:"updatedBy"`
}

type Certificate struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	LearnerID        string     `db:"learner_id" json:"learnerId"`
	SpecializationID *uuid.UUID `db:"specialization_id" json:"specializationId"`
	SubjectID        *uuid.UUID `db:"subject_id" json:"subjectId"`
	Name             string     `db:"name" json:"name"`
	Type             int32      `db:"type" json:"type"`
	Status           int32      `db:"status" json:"status"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
}

type Class struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Code             string    `db:"code" json:"code"`
	IsDraft          bool      `db:"is_draft" json:"isDraft"`
	Password         string    `db:"password" json:"password"`
	Name             string    `db:"name" json:"name"`
	Link             string    `db:"link" json:"link"`
	ProgramSubjectID uuid.UUID `db:"program_subject_id" json:"programSubjectId"`
	StartTime        time.Time `db:"start_time" json:"startTime"`
	EndTime          time.Time `db:"end_time" json:"endTime"`
	CreatedBy        string    `db:"created_by" json:"createdBy"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
}

type ClassLearner struct {
	ID        uuid.UUID `db:"id" json:"id"`
	LearnerID string    `db:"learner_id" json:"learnerId"`
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
}

type ClassTeacher struct {
	ID        uuid.UUID `db:"id" json:"id"`
	TeacherID string    `db:"teacher_id" json:"teacherId"`
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	CreatedBy string    `db:"created_by" json:"createdBy"`
}

type District struct {
	ID         int32  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	ProvinceID int32  `db:"province_id" json:"provinceId"`
}

type LearnerAssignment struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ClassLernerID uuid.UUID `db:"class_lerner_id" json:"classLernerId"`
	AssignmentID  uuid.UUID `db:"assignment_id" json:"assignmentId"`
	Grade         float64   `db:"grade" json:"grade"`
}

type LearnerAttendance struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ClassLearnerID uuid.UUID `db:"class_learner_id" json:"classLearnerId"`
	SlotID         uuid.UUID `db:"slot_id" json:"slotId"`
	Status         int32     `db:"status" json:"status"`
}

type LearnerSpecialization struct {
	LearnerID        string     `db:"learner_id" json:"learnerId"`
	SpecializationID uuid.UUID  `db:"specialization_id" json:"specializationId"`
	JoinedAt         *time.Time `db:"joined_at" json:"joinedAt"`
}

type Material struct {
	ID        uuid.UUID       `db:"id" json:"id"`
	SessionID uuid.UUID       `db:"session_id" json:"sessionId"`
	Index     int32           `db:"index" json:"index"`
	Type      string          `db:"type" json:"type"`
	Data      json.RawMessage `db:"data" json:"data"`
	IsShared  bool            `db:"is_shared" json:"isShared"`
	Name      *string         `db:"name" json:"name"`
}

type Program struct {
	ID          uuid.UUID   `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	StartDate   pgtype.Date `db:"start_date" json:"startDate"`
	EndDate     pgtype.Date `db:"end_date" json:"endDate"`
	CreatedBy   string      `db:"created_by" json:"createdBy"`
	UpdatedBy   *string     `db:"updated_by" json:"updatedBy"`
	Description string      `db:"description" json:"description"`
	CreatedAt   *time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time  `db:"updated_at" json:"updatedAt"`
}

type ProgramSubject struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	ProgramID uuid.UUID  `db:"program_id" json:"programId"`
	SubjectID uuid.UUID  `db:"subject_id" json:"subjectId"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
}

type Province struct {
	ID   int32  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type School struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Address    string    `db:"address" json:"address"`
	DistrictID int32     `db:"district_id" json:"districtId"`
	IsDeleted  *bool     `db:"is_deleted" json:"isDeleted"`
}

type Session struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SubjectID uuid.UUID `db:"subject_id" json:"subjectId"`
	Index     int32     `db:"index" json:"index"`
	Name      string    `db:"name" json:"name"`
}

type Skill struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

type Slot struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SessionID uuid.UUID `db:"session_id" json:"sessionId"`
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
	StartTime time.Time `db:"start_time" json:"startTime"`
	EndTime   time.Time `db:"end_time" json:"endTime"`
}

type Specialization struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Code        string     `db:"code" json:"code"`
	TimeAmount  *float64   `db:"time_amount" json:"timeAmount"`
	ImageLink   *string    `db:"image_link" json:"imageLink"`
	Status      int16      `db:"status" json:"status"`
	Description *string    `db:"description" json:"description"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	UpdatedBy   *string    `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
}

type SpecializationSkill struct {
	ID               uuid.UUID `db:"id" json:"id"`
	SpecializationID uuid.UUID `db:"specialization_id" json:"specializationId"`
	SkillID          uuid.UUID `db:"skill_id" json:"skillId"`
}

type SpecializationSubject struct {
	ID               uuid.UUID `db:"id" json:"id"`
	SpecializationID uuid.UUID `db:"specialization_id" json:"specializationId"`
	SubjectID        uuid.UUID `db:"subject_id" json:"subjectId"`
	CreatedBy        string    `db:"created_by" json:"createdBy"`
}

type Subject struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Code            string     `db:"code" json:"code"`
	Name            string     `db:"name" json:"name"`
	TimePerSession  int16      `db:"time_per_session" json:"timePerSession"`
	SessionsPerWeek int16      `db:"sessions_per_week" json:"sessionsPerWeek"`
	ImageLink       *string    `db:"image_link" json:"imageLink"`
	Status          int16      `db:"status" json:"status"`
	Description     string     `db:"description" json:"description"`
	CreatedBy       string     `db:"created_by" json:"createdBy"`
	UpdatedBy       *string    `db:"updated_by" json:"updatedBy"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
}

type SubjectSkill struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SubjectID uuid.UUID `db:"subject_id" json:"subjectId"`
	SkillID   uuid.UUID `db:"skill_id" json:"skillId"`
}

type Transcript struct {
	ID         uuid.UUID `db:"id" json:"id"`
	SubjectID  uuid.UUID `db:"subject_id" json:"subjectId"`
	Name       string    `db:"name" json:"name"`
	MinGrade   float64   `db:"min_grade" json:"minGrade"`
	MaxGrade   float64   `db:"max_grade" json:"maxGrade"`
	Percentage float64   `db:"percentage" json:"percentage"`
}

type User struct {
	ID           string     `db:"id" json:"id"`
	FullName     *string    `db:"full_name" json:"fullName"`
	Email        string     `db:"email" json:"email"`
	Phone        *string    `db:"phone" json:"phone"`
	Gender       *int16     `db:"gender" json:"gender"`
	AuthRole     int16      `db:"auth_role" json:"authRole"`
	ProfilePhoto *string    `db:"profile_photo" json:"profilePhoto"`
	Status       int32      `db:"status" json:"status"`
	SchoolID     *uuid.UUID `db:"school_id" json:"schoolId"`
}
