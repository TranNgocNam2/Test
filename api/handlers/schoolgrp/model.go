package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
)

func toCoreNewSchool(newSchoolRequest payload.NewSchool) school.NewSchool {
	return school.NewSchool{
		Name:       newSchoolRequest.Name,
		Address:    newSchoolRequest.Address,
		DistrictId: newSchoolRequest.DistrictId,
	}
}

func validateCreateSchoolRequest(newSchoolRequest payload.NewSchool) error {
	if err := validate.Check(newSchoolRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSchool(updateSchoolRequest payload.UpdateSchool) school.UpdateSchool {
	return school.UpdateSchool{
		Name:       updateSchoolRequest.Name,
		Address:    updateSchoolRequest.Address,
		DistrictId: updateSchoolRequest.DistrictId,
	}
}

func validateUpdateSchoolRequest(request payload.UpdateSchool) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
