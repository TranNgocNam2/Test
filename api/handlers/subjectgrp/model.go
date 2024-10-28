package subjectgrp

import (
	"Backend/internal/validate"
	"Backend/internal/web/payload"

	"gitlab.com/innovia69420/kit/web/request"
)

func validateNewSubjectRequest(request payload.NewSubject) error {
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
