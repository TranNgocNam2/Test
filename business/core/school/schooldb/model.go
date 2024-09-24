package schooldb

import (
	"Backend/business/core/school"
	"Backend/business/db/sqlc"
)

func toCoreSchool(dbSchool sqlc.School) school.School {
	return school.School{
		ID:      dbSchool.ID,
		Name:    dbSchool.Name,
		Address: dbSchool.Address,
	}
}

func toCoreSchoolSlice(dbSchools []sqlc.School) []school.School {
	schools := make([]school.School, len(dbSchools))
	for i, dbSchool := range dbSchools {
		schools[i] = toCoreSchool(dbSchool)
	}
	return schools
}

func toCoreProvince(dbProvince sqlc.Province) school.Province {
	return school.Province{
		ID:   dbProvince.ID,
		Name: dbProvince.Name,
	}
}

func toCoreProvinceSlice(dbProvinces []sqlc.Province) []school.Province {
	provinces := make([]school.Province, len(dbProvinces))
	for i, dbProvince := range dbProvinces {
		provinces[i] = toCoreProvince(dbProvince)
	}
	return provinces
}

func toCoreDistrict(dbDistrict sqlc.District) school.District {
	return school.District{
		ID:   dbDistrict.ID,
		Name: dbDistrict.Name,
	}
}

func toCoreDistrictSlice(dbDistricts []sqlc.District) []school.District {
	districts := make([]school.District, len(dbDistricts))
	for i, dbDistrict := range dbDistricts {
		districts[i] = toCoreDistrict(dbDistrict)
	}
	return districts
}
