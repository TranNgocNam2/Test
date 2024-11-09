package payload

type NewSpecialization struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required,uppercase"`
	TimeAmount  float64 `json:"timeAmount"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
}

type UpdateSpecialization struct {
	Name        string    `json:"name" validate:"required"`
	Code        string    `json:"code" validate:"required,uppercase"`
	Status      *int16    `json:"status" validate:"required,gte=0,lte=1"`
	TimeAmount  float64   `json:"timeAmount" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Subjects    []Subject `json:"subjects" validate:"required"`
}

type Subject struct {
	ID    string `json:"id" validate:"required"`
	Index *int   `json:"index" validate:"required,gte=0"`
}
