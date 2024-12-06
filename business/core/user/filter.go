package user

import (
	"Backend/internal/validate"
	"bytes"
	"fmt"
	"strings"
)

const (
	UserStatus     = "u.status"
	VerifiedStatus = "vl.status"
)

type QueryFilter struct {
	FullName   *string `validate:"omitempty"`
	Status     *int16  `validate:"omitempty"`
	SchoolName *string `validate:"omitempty"`
	Role       *int    `validate:"omitempty"`
	IsVerified *bool   `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return err
	}
	return nil
}

func (qf *QueryFilter) WithIsVerified(isVerified bool) {
	qf.IsVerified = &isVerified
}

func (qf *QueryFilter) WithFullName(fullName string) {
	qf.FullName = &fullName
}

func (qf *QueryFilter) WithStatus(status int16) {
	qf.Status = &status
}

func (qf *QueryFilter) WithSchoolName(schoolName string) {
	qf.SchoolName = &schoolName
}

func (qf *QueryFilter) WithRole(role int) {
	qf.Role = &role
}

func applyFilter(filter QueryFilter, data map[string]interface{}, buf *bytes.Buffer, hasWhere bool, statusType string) {
	var wc []string

	if filter.FullName != nil {
		data["full_name"] = fmt.Sprintf("%%%s%%", *filter.FullName)
		wc = append(wc, "u.full_name LIKE :full_name")
	}

	if filter.SchoolName != nil {
		data["school_name"] = fmt.Sprintf("%%%s%%", *filter.SchoolName)
		wc = append(wc, "s.name LIKE :school_name")
	}

	if filter.IsVerified != nil {
		data["is_verified"] = fmt.Sprintf("%t", *filter.IsVerified)
		wc = append(wc, "u.is_verified = :is_verified")
	}

	if filter.Status != nil {
		if statusType == UserStatus {
			data["status"] = fmt.Sprintf("%d", *filter.Status)
			wc = append(wc, "u.status = :status")
		}
		if statusType == VerifiedStatus {
			data["status"] = fmt.Sprintf("%d", *filter.Status)
			wc = append(wc, "vl.status = :status")
		}
	}

	if filter.Role != nil {
		data["role"] = fmt.Sprintf("%d", *filter.Role)
		wc = append(wc, "u.auth_role = :role")
	}

	if len(wc) > 0 {
		if !hasWhere {
			buf.WriteString(" WHERE ")
			hasWhere = true
		} else {
			buf.WriteString(" AND ")
		}

		buf.WriteString(strings.Join(wc, " AND "))
	}
}
