package transcript

import (
	"bytes"
	"fmt"
	"strings"
)

type QueryFilter struct {
	TranscriptName *string `validate:"omitempty"`
}

func (qf *QueryFilter) WithName(name string) {
	qf.TranscriptName = &name
}

func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.TranscriptName != nil {
		data["transcript_name"] = fmt.Sprintf("%%%s%%", *filter.TranscriptName)
		wc = append(wc, "t.name LIKE :transcript_name")
	}

	if len(wc) > 0 {
		buf.WriteString(" AND ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
