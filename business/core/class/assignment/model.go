package assignment

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	Id         uuid.UUID       `json:"id"`
	ClassId    uuid.UUID       `json:"classId"`
	Question   json.RawMessage `json:"question"`
	Deadline   time.Time       `json:"deadline"`
	Status     int             `json:"status"`
	Type       int             `json:"type"`
	CanOverdue bool            `json:"canOverdue"`
}
