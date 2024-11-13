package learner

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
)

type ClassAccess struct {
	Code     string
	Password string
}

type AttendanceSubmission struct {
	Index          int32
	AttendanceCode string
}

type Learner struct {
	ID          string       `json:"id"`
	FullName    string       `json:"fullName"`
	Email       string       `json:"email"`
	Phone       string       `json:"phone"`
	Gender      *int16       `json:"gender"`
	Photo       string       `json:"photo"`
	School      School       `json:"school"`
	ImageLink   []string     `json:"image_links"`
	Attendances []Attendance `json:"attendances"`
	Assignments []Assignment `json:"assignments"`
}

type School struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Attendance struct {
	ID     uuid.UUID `json:"id"`
	SlotId uuid.UUID `json:"slotId"`
	Status int32     `json:"status"`
}

type AttendanceRecord struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	School   School `json:"school"`
	Status   int32  `json:"status"`
}

func toCoreAttendance(dbAttendance sqlc.LearnerAttendance) Attendance {
	return Attendance{
		ID:     dbAttendance.ID,
		SlotId: dbAttendance.SlotID,
		Status: dbAttendance.Status,
	}
}

func toCoreAttendanceSlice(dbAttendances []sqlc.LearnerAttendance) []Attendance {
	attendances := make([]Attendance, len(dbAttendances))
	for i, dbAttendance := range dbAttendances {
		attendances[i] = toCoreAttendance(dbAttendance)
	}
	return attendances
}

type Assignment struct {
	ID           uuid.UUID `json:"id"`
	AssignmentId uuid.UUID `json:"assignmentId"`
	Grade        float32   `json:"grade"`
}

func toCoreAssignment(dbLearnerAssignment sqlc.LearnerAssignment) Assignment {
	return Assignment{
		ID:           dbLearnerAssignment.ID,
		AssignmentId: dbLearnerAssignment.AssignmentID,
		Grade:        dbLearnerAssignment.Grade,
	}
}

func toCoreAssignmentSlice(dbLearnerAssignments []sqlc.LearnerAssignment) []Assignment {
	assignments := make([]Assignment, len(dbLearnerAssignments))
	for i, dbLearnerAssignment := range dbLearnerAssignments {
		assignments[i] = toCoreAssignment(dbLearnerAssignment)
	}
	return assignments
}
