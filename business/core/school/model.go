package school

import "github.com/google/uuid"

type School struct {
	ID         uuid.UUID `db:"id"`
	Name       string    `db:"name"`
	Address    string    `db:"address"`
	DistrictID int32     `db:"district_id"`
}

type NewSchool struct {
	Name       string
	Address    string
	DistrictID int32
}

type UpdateSchool struct {
	Name       *string
	Address    *string
	DistrictID *int32
}

type Province struct {
	ID   int32
	Name string
}

type District struct {
	ID         int32
	Name       string
	ProvinceID int32
}
