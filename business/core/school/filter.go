package school

import (
	"Backend/internal/validate"
	"bytes"
	"fmt"
	"strings"
)

type QueryFilter struct {
	Name *string `validate:"omitempty"`
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

func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "s.name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
		buf.WriteString(" AND s.is_deleted = false")
	} else {
		buf.WriteString(" WHERE s.is_deleted = false")
	}
}
