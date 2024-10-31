package subjectgrp

import (
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"encoding/json"
	"fmt"
)

func invalidMaterialError(material payload.Material) error {
	return fmt.Errorf("Data cho material với id %s không đúng định dạng!", material.ID)
}

type CodeData struct {
	Value    string `json:"value" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type FileData struct {
	Source   string `json:"src" validate:"required"`
	FileName string `json:"fileName" validate:"required"`
	FileSize int    `json:"fileSize" validate:"required"`
}

func validateMaterial(material payload.Material) error {
	switch material.Type {
	case "h1", "h2", "h3", "text", "image", "video":
		if _, ok := material.Data.(string); !ok {
			return validate.NewFieldsError("materials", invalidMaterialError(material))
		}
		return nil
	case "file":
		data, err := json.Marshal(material.Data)
		if err != nil {
			return validate.NewFieldsError("materials", invalidMaterialError(material))
		}

		var materialType FileData
		if err := json.Unmarshal(data, &materialType); err != nil {
			return validate.NewFieldsError("materials", invalidMaterialError(material))
		}

		if err := validate.Check(materialType); err != nil {
			return err
		}

		return nil
	case "code":
		data, err := json.Marshal(material.Data)
		if err != nil {
			return validate.NewFieldsError("materials", invalidMaterialError(material))
		}

		var materialType CodeData
		if err := json.Unmarshal(data, &materialType); err != nil {
			return validate.NewFieldsError("materials", invalidMaterialError(material))
		}

		if err := validate.Check(materialType); err != nil {
			return err
		}

		return nil
	default:
		return validate.NewFieldsError("materials", invalidMaterialError(material))
	}
}

func validateNewSubjectRequest(request payload.NewSubject) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}

func validateUpdateSubjectRequest(request payload.UpdateSubject) error {
	if err := validate.Check(request); err != nil {
		return err
	}

	// update published request will have no sessions
	if len(request.Sessions) == 0 {
		return nil
	}

	for _, session := range request.Sessions {
		if len(session.Materials) == 0 {
			continue
		}

		for _, material := range session.Materials {
			if err := validateMaterial(material); err != nil {
				return err
			}
		}
	}
	return nil
}
