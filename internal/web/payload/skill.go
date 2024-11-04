package payload

type NewSkill struct {
	Name string `json:"name" validate:"required"`
}

type UpdateSkill struct {
	Name string `json:"name" validate:"required"`
}
