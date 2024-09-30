package schoolgrp

import (
	"Backend/business/core/school"
	"github.com/google/uuid"
)

type SchoolResponse struct {
	ID         uuid.UUID `json:"id"`
	SchoolName string    `json:"schoolName"`
	Address    string    `json:"address"`
}

func toSchoolResponse(school school.School) SchoolResponse {
	return SchoolResponse{
		ID:         school.ID,
		SchoolName: school.Name,
		Address:    school.Address,
	}
}

func toWebSchools(schools []school.School) []SchoolResponse {
	items := make([]SchoolResponse, len(schools))
	for i, school := range schools {
		items[i] = toSchoolResponse(school)
	}
	return items
}

type ProvinceResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func toProvinceResponse(province school.Province) ProvinceResponse {
	return ProvinceResponse{
		ID:   province.ID,
		Name: province.Name,
	}
}
func toProvinceResponses(provinces []school.Province) []ProvinceResponse {
	items := make([]ProvinceResponse, len(provinces))
	for i, dbProvince := range provinces {
		items[i] = toProvinceResponse(dbProvince)
	}
	return items
}

type DistrictResponse struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	ProvinceID int32  `json:"provinceID"`
}

func toDistrictResponse(district school.District) DistrictResponse {
	return DistrictResponse{
		ID:         district.ID,
		Name:       district.Name,
		ProvinceID: district.ProvinceID,
	}
}
func toClientDistricts(districts []school.District) []DistrictResponse {
	items := make([]DistrictResponse, len(districts))
	for i, district := range districts {
		items[i] = toDistrictResponse(district)
	}
	return items
}
