package transcript

import (
	"Backend/internal/validate"
	"bytes"
	"fmt"
	"strings"
)

type QueryFilter struct {
	TranscriptName *string `validate:"omitempty"`
	LearnerId      *string `validate:"omitempty"`
}

func (qf *QueryFilter) WithName(name string) {
	qf.TranscriptName = &name
}

func (qf *QueryFilter) WithLearnerId(id string) {
	qf.LearnerId = &id
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return err
	}

	return nil
}

func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.TranscriptName != nil {
		data["transcript_name"] = fmt.Sprintf("%%%s%%", *filter.TranscriptName)
		wc = append(wc, "t.name LIKE :transcript_name")
	}

	if filter.LearnerId != nil {
		data["learner_id"] = *filter.LearnerId
		wc = append(wc, "cl.learner_id = :learner_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" AND ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
