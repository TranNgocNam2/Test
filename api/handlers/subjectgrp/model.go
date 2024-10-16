package subjectgrp

import (
	"Backend/internal/validate"

	"gitlab.com/innovia69420/kit/web/request"
)

func validateNewSubjectRequest(request request.NewSubject) error {
	if err := validate.Check(request); err != nil {
		return err
	}

	return nil
}

func validateUpdateSubjectRequest(request request.UpdateSubject) error {
	if err := validate.Check(request); err != nil {
		return err
	}

	return nil
}
