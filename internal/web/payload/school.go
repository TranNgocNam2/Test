package payload

type NewSchool struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	DistrictID int32  `json:"districtID" validate:"required"`
}

type UpdateSchool struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	DistrictID int32  `json:"districtID" validate:"required"`
}
