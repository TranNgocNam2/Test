package payload

type UpdateSlot struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	TeacherId string `json:"teacherId"`
}
