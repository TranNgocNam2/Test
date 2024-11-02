package payload

type NewProgram struct {
	Name        string `json:"name" validate:"required"`
	StartDate   string `json:"startDate" validate:"required"`
	EndDate     string `json:"endDate" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateProgram struct {
	Name        string `json:"name" validate:"required"`
	StartDate   string `json:"startDate" validate:"required"`
	EndDate     string `json:"endDate" validate:"required"`
	Description string `json:"description" validate:"required"`
}
