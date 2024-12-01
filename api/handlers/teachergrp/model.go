package teachergrp

import (
	"Backend/business/core/teacher"
	"Backend/internal/validate"
	"Backend/internal/web/payload"
)

func toCoreUpdateRecord(request payload.UpdateRecord) teacher.UpdateRecord {
	return teacher.UpdateRecord{
		Link: request.Link,
	}
}

func validateUpdateRecordRequest(request payload.UpdateRecord) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
