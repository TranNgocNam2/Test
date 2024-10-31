package payload

type NewSubject struct {
	Name           string   `json:"name" validate:"required"`
	Code           string   `json:"code" validate:"required"`
	Description    string   `json:"description"`
	Image          string   `json:"image" validate:"required"`
	TimePerSession int      `json:"timePerSession" validate:"required"`
	Skills         []string `json:"skills" validate:"gt=0,dive,required"`
}

type UpdateSubject struct {
	Name           string       `json:"name" validate:"required"`
	Code           string       `json:"code" validate:"required"`
	Image          string       `json:"image" validate:"required"`
	TimePerSession int          `json:"timePerSession" validate:"required"`
	MinPassGrade   float64      `json:"minPassGrade" validate:"required"`
	MinAttendance  float64      `json:"minAttendance" validate:"required"`
	Description    string       `json:"description"`
	Status         *int         `json:"status" validate:"gte=0,lte=1,required"`
	Skills         []string     `json:"skills" validate:"gt=0,dive,required"`
	Sessions       []Session    `json:"sessions"`
	Transcripts    []Transcript `json:"transcripts"`
}

type Transcript struct {
	Id         string  `json:"id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Index      int     `json:"index" validate:"required"`
	Percentage float32 `json:"percentage" validate:"required"`
	MinGrade   float32 `json:"minGrade" validate:"required"`
}

type Session struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Index     int        `json:"index"`
	Materials []Material `json:"materials"`
}

type Material struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Index    int         `json:"index"`
	IsShared bool        `json:"isShared"`
	Data     interface{} `json:"data"`
}
