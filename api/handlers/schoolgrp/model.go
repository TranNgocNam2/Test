package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/validate"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
)

type SchoolResponse struct {
	ID         uuid.UUID `json:"id"`
	SchoolName string    `json:"schoolName"`
	Address    string    `json:"address"`
	DistrictID int       `json:"districtID"`
}

func toSchoolResponse(school school.School) SchoolResponse {
	return SchoolResponse{
		ID:         school.ID,
		SchoolName: school.Name,
		Address:    school.Address,
		DistrictID: int(school.DistrictID),
	}
}

func toSchoolsResponse(schools []school.School) []SchoolResponse {
	items := make([]SchoolResponse, len(schools))
	for i, school := range schools {
		items[i] = toSchoolResponse(school)
	}
	return items
}

func toCoreNewSchool(newSchoolRequest request.NewSchool) school.School {
	return school.School{
		Name:       newSchoolRequest.Name,
		Address:    newSchoolRequest.Address,
		DistrictID: newSchoolRequest.DistrictID,
	}
}

func validateCreateSchoolRequest(newSchoolRequest request.NewSchool) error {
	if err := validate.Check(newSchoolRequest); err != nil {
		return err
	}
	return nil
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

func toDistrictsResponse(districts []school.District) []DistrictResponse {
	items := make([]DistrictResponse, len(districts))
	for i, district := range districts {
		items[i] = toDistrictResponse(district)
	}
	return items
}

func toCoreUpdateSchool(updateSchoolRequest request.UpdateSchool) school.School {
	return school.School{
		Name:       updateSchoolRequest.Name,
		Address:    updateSchoolRequest.Address,
		DistrictID: updateSchoolRequest.DistrictID,
	}
}

func validateUpdateSchoolRequest(request request.UpdateSchool) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
