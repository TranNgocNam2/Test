package payload

import "encoding/json"

type NewSubject struct {
	Name           string   `json:"name" validate:"required"`
	Code           string   `json:"code" validate:"required"`
	Description    string   `json:"description"`
	Image          string   `json:"image" validate:"required"`
	TimePerSession int      `json:"timePerSession" validate:"required"`
	Skills         []string `json:"skills" validate:"gt=0,dive,required"`
}

type UpdateSubject struct {
	Name        string    `json:"name" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Description string    `json:"description"`
	Status      *int      `json:"status" validate:"gte=0,lte=1,required"`
	Skills      []string  `json:"skills" validate:"gt=0,dive,required"`
	Sessions    []Session `json:"sessions"`
}

type Session struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Index     int        `json:"index"`
	Materials []Material `json:"materials"`
}

type Material struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Index    int             `json:"index"`
	IsShared bool            `json:"isShared"`
	Data     json.RawMessage `json:"data"`
}
