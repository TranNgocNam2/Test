package school

import (
	"Backend/business/db/sqlc"

	"github.com/google/uuid"
)

// School

type School struct {
	ID         uuid.UUID
	Name       string
	Address    string
	DistrictID int32
}

func toCoreSchool(dbSchool sqlc.School) School {
	return School{
		ID:         dbSchool.ID,
		Name:       dbSchool.Name,
		Address:    dbSchool.Address,
		DistrictID: dbSchool.DistrictID,
	}
}

func toCoreSchoolSlice(dbSchools []sqlc.School) []School {
	schools := make([]School, len(dbSchools))
	for i, school := range dbSchools {
		schools[i] = toCoreSchool(school)
	}

	return schools
}

// Province

type Province struct {
	ID   int32
	Name string
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
	ID         int32
	Name       string
	ProvinceID int32
}

func toCoreDistrict(dbDistrict sqlc.District) District {
	return District{
		ID:         dbDistrict.ID,
		Name:       dbDistrict.Name,
		ProvinceID: dbDistrict.ProvinceID,
	}
}

func toCoreDistrictSlice(dbDistricts []sqlc.District) []District {
	districts := make([]District, len(dbDistricts))
	for i, dbDistrict := range dbDistricts {
		districts[i] = toCoreDistrict(dbDistrict)
	}
	return districts
}
