package schooldb

import (
	"Backend/business/core/school"
	"bytes"
	"fmt"
	"strings"
)

func applyFilter(filter school.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
