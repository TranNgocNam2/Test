package schoolgrp

import (
	"Backend/business/db/sqlc"
	"Backend/internal/validate"
	"fmt"
	"github.com/google/uuid"
)

type ClientNewSchool struct {
	SchoolName string `json:"schoolName" validate:"required"`
	Address    string `json:"address" validate:"required"`
	DistrictID int32  `json:"districtID" validate:"required"`
}

func toCoreNewSchool(clientNewSchool ClientNewSchool) sqlc.CreateSchoolParams {
	return sqlc.CreateSchoolParams{
		ID:         uuid.New(),
		Name:       clientNewSchool.SchoolName,
		Address:    clientNewSchool.Address,
		DistrictID: clientNewSchool.DistrictID,
	}
}

func (clientNewSchool ClientNewSchool) Validate() error {
	if err := validate.Check(clientNewSchool); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

type ClientProvince struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type ClientDistrict struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	ProvinceID int32  `json:"province_id"`
}

func toClientProvince(dbProvince sqlc.Province) ClientProvince {
	return ClientProvince{
		ID:   dbProvince.ID,
		Name: dbProvince.Name,
	}
}
func toClientProvinces(dbProvinces []sqlc.Province) []ClientProvince {
	provinces := make([]ClientProvince, len(dbProvinces))
	for i, dbProvince := range dbProvinces {
		provinces[i] = toClientProvince(dbProvince)
	}
	return provinces
}

func toClientDistrict(dbDistrict sqlc.District) ClientDistrict {
	return ClientDistrict{
		ID:         dbDistrict.ID,
		Name:       dbDistrict.Name,
		ProvinceID: dbDistrict.ProvinceID,
	}
}
func toClientDistricts(dbDistricts []sqlc.District) []ClientDistrict {
	districts := make([]ClientDistrict, len(dbDistricts))
	for i, dbDistrict := range dbDistricts {
		districts[i] = toClientDistrict(dbDistrict)
	}
	return districts
}
