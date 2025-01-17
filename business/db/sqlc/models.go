// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	ClassID    uuid.UUID       `db:"class_id" json:"classId"`
	Question   json.RawMessage `db:"question" json:"question"`
	Deadline   *time.Time      `db:"deadline" json:"deadline"`
	Status     int16           `db:"status" json:"status"`
	CanOverdue *bool           `db:"can_overdue" json:"canOverdue"`
	Type       int16           `db:"type" json:"type"`
}

type Certificate struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	LearnerID        string     `db:"learner_id" json:"learnerId"`
	SpecializationID *uuid.UUID `db:"specialization_id" json:"specializationId"`
	SubjectID        *uuid.UUID `db:"subject_id" json:"subjectId"`
	ClassID          *uuid.UUID `db:"class_id" json:"classId"`
	Name             string     `db:"name" json:"name"`
	Status           int32      `db:"status" json:"status"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy        *string    `db:"updated_by" json:"updatedBy"`
}

type Class struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Code      string     `db:"code" json:"code"`
	SubjectID uuid.UUID  `db:"subject_id" json:"subjectId"`
	ProgramID uuid.UUID  `db:"program_id" json:"programId"`
	Password  string     `db:"password" json:"password"`
	Name      string     `db:"name" json:"name"`
	Link      *string    `db:"link" json:"link"`
	StartDate *time.Time `db:"start_date" json:"startDate"`
	EndDate   *time.Time `db:"end_date" json:"endDate"`
	Status    int16      `db:"status" json:"status"`
	Type      int16      `db:"type" json:"type"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *string    `db:"updated_by" json:"updatedBy"`
}

type ClassLearner struct {
	ID        uuid.UUID `db:"id" json:"id"`
	LearnerID string    `db:"learner_id" json:"learnerId"`
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
	Status    int16     `db:"status" json:"status"`
}

type District struct {
	ID         int32  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	ProvinceID int32  `db:"province_id" json:"provinceId"`
}

type LearnerAssignment struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ClassLearnerID   uuid.UUID `db:"class_learner_id" json:"classLearnerId"`
	AssignmentID     uuid.UUID `db:"assignment_id" json:"assignmentId"`
	Grade            float32   `db:"grade" json:"grade"`
	Data             []byte    `db:"data" json:"data"`
	GradingStatus    int16     `db:"grading_status" json:"gradingStatus"`
	SubmissionStatus int16     `db:"submission_status" json:"submissionStatus"`
}

type LearnerAttendance struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ClassLearnerID uuid.UUID `db:"class_learner_id" json:"classLearnerId"`
	SlotID         uuid.UUID `db:"slot_id" json:"slotId"`
	Status         int32     `db:"status" json:"status"`
}

type LearnerSpecialization struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	LearnerID        string     `db:"learner_id" json:"learnerId"`
	SpecializationID uuid.UUID  `db:"specialization_id" json:"specializationId"`
	JoinedAt         *time.Time `db:"joined_at" json:"joinedAt"`
}

type LearnerTranscript struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ClassLearnerID uuid.UUID  `db:"class_learner_id" json:"classLearnerId"`
	TranscriptID   uuid.UUID  `db:"transcript_id" json:"transcriptId"`
	Grade          *float32   `db:"grade" json:"grade"`
	Status         int16      `db:"status" json:"status"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy      *string    `db:"updated_by" json:"updatedBy"`
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
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	StartDate   time.Time  `db:"start_date" json:"startDate"`
	EndDate     time.Time  `db:"end_date" json:"endDate"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	UpdatedBy   *string    `db:"updated_by" json:"updatedBy"`
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
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
	ID             uuid.UUID  `db:"id" json:"id"`
	SessionID      uuid.UUID  `db:"session_id" json:"sessionId"`
	ClassID        uuid.UUID  `db:"class_id" json:"classId"`
	StartTime      *time.Time `db:"start_time" json:"startTime"`
	EndTime        *time.Time `db:"end_time" json:"endTime"`
	Index          int32      `db:"index" json:"index"`
	TeacherID      *string    `db:"teacher_id" json:"teacherId"`
	AttendanceCode *string    `db:"attendance_code" json:"attendanceCode"`
	RecordLink     *string    `db:"record_link" json:"recordLink"`
}

type Specialization struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Code        string     `db:"code" json:"code"`
	TimeAmount  *float32   `db:"time_amount" json:"timeAmount"`
	ImageLink   *string    `db:"image_link" json:"imageLink"`
	Status      int16      `db:"status" json:"status"`
	Description *string    `db:"description" json:"description"`
	CreatedBy   string     `db:"created_by" json:"createdBy"`
	UpdatedBy   *string    `db:"updated_by" json:"updatedBy"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
}

type SpecializationSubject struct {
	ID               uuid.UUID `db:"id" json:"id"`
	SpecializationID uuid.UUID `db:"specialization_id" json:"specializationId"`
	SubjectID        uuid.UUID `db:"subject_id" json:"subjectId"`
	Index            int16     `db:"index" json:"index"`
	CreatedBy        string    `db:"created_by" json:"createdBy"`
}

type Subject struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Code            string     `db:"code" json:"code"`
	Name            string     `db:"name" json:"name"`
	TimePerSession  float32    `db:"time_per_session" json:"timePerSession"`
	SessionsPerWeek int16      `db:"sessions_per_week" json:"sessionsPerWeek"`
	TotalSessions   int16      `db:"total_sessions" json:"totalSessions"`
	MinPassGrade    *float32   `db:"min_pass_grade" json:"minPassGrade"`
	MinAttendance   *float32   `db:"min_attendance" json:"minAttendance"`
	ImageLink       *string    `db:"image_link" json:"imageLink"`
	Status          int16      `db:"status" json:"status"`
	Description     *string    `db:"description" json:"description"`
	CreatedBy       string     `db:"created_by" json:"createdBy"`
	UpdatedBy       *string    `db:"updated_by" json:"updatedBy"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	LearnerType     *int16     `db:"learner_type" json:"learnerType"`
}

type SubjectSkill struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SubjectID uuid.UUID `db:"subject_id" json:"subjectId"`
	SkillID   uuid.UUID `db:"skill_id" json:"skillId"`
}

type Transcript struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SubjectID uuid.UUID `db:"subject_id" json:"subjectId"`
	Name      string    `db:"name" json:"name"`
	Index     int32     `db:"index" json:"index"`
	MinGrade  float64   `db:"min_grade" json:"minGrade"`
	Weight    float64   `db:"weight" json:"weight"`
}

type User struct {
	ID           string     `db:"id" json:"id"`
	FullName     *string    `db:"full_name" json:"fullName"`
	Email        string     `db:"email" json:"email"`
	Phone        *string    `db:"phone" json:"phone"`
	AuthRole     int16      `db:"auth_role" json:"authRole"`
	ProfilePhoto *string    `db:"profile_photo" json:"profilePhoto"`
	Status       int32      `db:"status" json:"status"`
	IsVerified   bool       `db:"is_verified" json:"isVerified"`
	SchoolID     *uuid.UUID `db:"school_id" json:"schoolId"`
	Type         *int16     `db:"type" json:"type"`
}

type VerificationLearner struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	SchoolID   uuid.UUID  `db:"school_id" json:"schoolId"`
	LearnerID  string     `db:"learner_id" json:"learnerId"`
	ImageLink  []string   `db:"image_link" json:"imageLink"`
	Status     int16      `db:"status" json:"status"`
	VerifiedBy *string    `db:"verified_by" json:"verifiedBy"`
	Type       int16      `db:"type" json:"type"`
	VerifiedAt *time.Time `db:"verified_at" json:"verifiedAt"`
	Note       *string    `db:"note" json:"note"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
}
