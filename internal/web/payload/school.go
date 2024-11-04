package payload

type NewSchool struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	DistrictId int32  `json:"districtId" validate:"required"`
}

type UpdateSchool struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	DistrictId int32  `json:"districtId" validate:"required"`
}
