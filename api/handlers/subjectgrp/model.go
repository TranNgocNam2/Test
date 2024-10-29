package subjectgrp

import (
	"Backend/internal/validate"
	"Backend/internal/web/payload"
	"encoding/json"
	"errors"
	"fmt"
)

func invalidMaterialError(material payload.Material) error {
	return fmt.Errorf("Data cho material với id %s không đúng định dạng!", material.ID)
}

var ErrSussyBaka = errors.New("Ko the nao!")

type CodeData struct {
	Value    string `json:"value" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type FileData struct {
	Source   string `json:"src" validate:"required"`
	FileName string `json:"fileName" validate:"required"`
	FileSize int    `json:"fileSize" validate:"required"`
}

func validateMaterials(materials []payload.Material) error {
	for _, material := range materials {
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
	return ErrSussyBaka
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
		if err := validateMaterials(session.Materials); err != nil {
			return err
		}
	}
	return nil
}
