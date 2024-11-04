package skillgrp

import (
	"Backend/business/core/skill"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"github.com/google/uuid"
)

func toCoreNewSkill(newSkill payload.NewSkill) skill.NewSkill {
	return skill.NewSkill{
		ID:   uuid.New(),
		Name: newSkill.Name,
	}
}

func validateNewSkillRequest(newSpecializationRequest payload.NewSkill) error {
	if err := validate.Check(newSpecializationRequest); err != nil {
		return err
	}
	return nil
}

func toCoreUpdateSkill(updateSkill payload.UpdateSkill) skill.UpdateSkill {
	return skill.UpdateSkill{
		Name: updateSkill.Name,
	}
}

func validateUpdateSkillRequest(updateSpecializationRequest payload.UpdateSkill) error {
	if err := validate.Check(updateSpecializationRequest); err != nil {
		return err
	}
	return nil
}
