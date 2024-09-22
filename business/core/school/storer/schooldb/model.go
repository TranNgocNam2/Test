package schooldb

import (
	"Backend/business/core/school"
	"Backend/business/db/sqlc"
)

//func toDBSchool(school school.NewSchool) sqlc.School {
//	return sqlc.School{
//		ID:        school.ID,
//		Name:      school.Name,
//		Address:   school.Address,
//		CreatedBy: school.CreatedBy,
//	}
//}

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
