package specialization

import (
	"Backend/internal/validate"
	"bytes"
	"fmt"
	"strings"
)

type QueryFilter struct {
	Name   *string `validate:"omitempty"`
	Code   *string `validate:"omitempty"`
	Status int16   `validate:"omitempty"`
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

func (qf *QueryFilter) WithCode(code string) {
	qf.Code = &code
}

func (qf *QueryFilter) WithStatus(status int16) {
	qf.Status = status
}

func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.Code != nil {
		data["code"] = *filter.Code
		wc = append(wc, "code = :code")
	}

	data["status"] = fmt.Sprintf("%d", filter.Status)
	wc = append(wc, "status = :status")

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
