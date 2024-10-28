// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: class_teachers.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const getTeachersByClassID = `-- name: GetTeachersByClassID :many
SELECT t.id, t.full_name, t.email, t.phone, t.gender, t.auth_role, t.profile_photo, t.status, t.school_id
FROM class_teachers ct
JOIN users t ON ct.teacher_id = t.id
WHERE ct.class_id = $1::uuid
`

func (q *Queries) GetTeachersByClassID(ctx context.Context, classID uuid.UUID) ([]User, error) {
	rows, err := q.db.Query(ctx, getTeachersByClassID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.Gender,
			&i.AuthRole,
			&i.ProfilePhoto,
			&i.Status,
			&i.SchoolID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
