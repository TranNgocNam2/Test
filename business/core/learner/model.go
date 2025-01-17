package learner

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type ClassAccess struct {
	Code     string
	Password string
}

type AttendanceSubmission struct {
	AttendanceCode string
}

type Learner struct {
	ID          string       `json:"id"`
	FullName    string       `json:"fullName"`
	Email       string       `json:"email"`
	Phone       *string      `json:"phone"`
	Photo       *string      `json:"photo"`
	Type        *int16       `json:"type"`
	School      School       `json:"school"`
	Attendances []Attendance `json:"attendances"`
	//Assignments []Assignment `json:"assignments"`
}

type VerifyLearnerInfo struct {
	ID            string `json:"id"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	Verifications []struct {
		ID        uuid.UUID `json:"id"`
		Status    int16     `json:"status"`
		Note      *string   `json:"note"`
		ImageLink []string  `json:"imageLink"`
		Type      int16     `json:"type"`
		School    School    `json:"school"`
		CreatedAt time.Time `json:"createdAt"`
	} `json:"verifications"`
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

type AttendanceReport struct {
	ID        uuid.UUID `json:"id"`
	Index     int32     `json:"index"`
	Status    int32     `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type UpdateLearner struct {
	SchoolId   uuid.UUID
	ImageLinks []string
	Type       int16
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

//type Assignment struct {
//	ID           uuid.UUID `json:"id"`
//	AssignmentId uuid.UUID `json:"assignmentId"`
//	Grade        float32   `json:"grade"`
//}

//func toCoreAssignment(dbLearnerAssignment sqlc.LearnerAssignment) Assignment {
//	return Assignment{
//		ID:           dbLearnerAssignment.ID,
//		AssignmentId: dbLearnerAssignment.AssignmentID,
//		Grade:        dbLearnerAssignment.Grade,
//	}
//}
//
//func toCoreAssignmentSlice(dbLearnerAssignments []sqlc.LearnerAssignment) []Assignment {
//	assignments := make([]Assignment, len(dbLearnerAssignments))
//	for i, dbLearnerAssignment := range dbLearnerAssignments {
//		assignments[i] = toCoreAssignment(dbLearnerAssignment)
//	}
//	return assignments
//}
