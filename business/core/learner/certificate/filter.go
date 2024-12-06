package certificate

import (
	"Backend/internal/validate"
	"bytes"
	"fmt"
	"strings"
)

type QueryFilter struct {
	Name      *string `validate:"omitempty"`
	LearnerId *string `validate:"omitempty"`
	Status    int16   `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return err
	}

	return nil
}

func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

func (qf *QueryFilter) WithLearnerId(learnerId string) {
	qf.LearnerId = &learnerId
}

func (qf *QueryFilter) WithStatus(status int16) {
	qf.Status = status
}
func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer, hasWhere bool) {
	var wc []string

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.LearnerId != nil {
		data["learnerId"] = *filter.LearnerId
		wc = append(wc, "c.learner_id = :learnerId")
	}

	data["status"] = filter.Status
	wc = append(wc, "c.status = :status")

	if len(wc) > 0 {
		if !hasWhere {
			buf.WriteString(" WHERE ")
		} else {
			buf.WriteString(" AND ")
		}
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
