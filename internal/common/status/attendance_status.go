package status

type Attendance int16

const (
	NotStarted Attendance = iota
	Attended
	Absent
)
