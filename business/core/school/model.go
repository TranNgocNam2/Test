package school

import (
	"github.com/google/uuid"
)

type School struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

type Province struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type District struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProvinceID int    `json:"province_id"`
}

//type NewSchool struct {
//	ID        uuid.UUID `json:"id"`
//	Name      string    `json:"name"`
//	Address   string    `json:"address"`
//	CreatedBy uuid.UUID `json:"created_by"`
//}
