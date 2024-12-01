package school

import (
	"Backend/business/db/sqlc"

	"github.com/google/uuid"
)

type NewSchool struct {
	Name       string
	Address    string
	DistrictId int32
}

type UpdateSchool struct {
	Name       string
	Address    string
	DistrictId int32
}

// School

type School struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	District District  `json:"district"`
	Province Province  `json:"province"`
}

func toCoreSchool(dbSchool sqlc.School) School {
	return School{
		ID:      dbSchool.ID,
		Name:    dbSchool.Name,
		Address: dbSchool.Address,
	}
}

// Province
type Province struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func toCoreProvince(dbProvince sqlc.Province) Province {
	return Province{
		ID:   dbProvince.ID,
		Name: dbProvince.Name,
	}
}

func toCoreProvinceSlice(dbProvinces []sqlc.Province) []Province {
	provinces := make([]Province, len(dbProvinces))
	for i, dbProvince := range dbProvinces {
		provinces[i] = toCoreProvince(dbProvince)
	}
	return provinces
}

// District

type District struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func toCoreDistrict(dbDistrict sqlc.District) District {
	return District{
		ID:   dbDistrict.ID,
		Name: dbDistrict.Name,
	}
}

func toCoreDistrictSlice(dbDistricts []sqlc.District) []District {
	districts := make([]District, len(dbDistricts))
	for i, dbDistrict := range dbDistricts {
		districts[i] = toCoreDistrict(dbDistrict)
	}
	return districts
}
