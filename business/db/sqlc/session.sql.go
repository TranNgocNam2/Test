// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: session.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const countSessionsBySubjectID = `-- name: CountSessionsBySubjectID :one
SELECT count(*) FROM sessions WHERE subject_id = $1
`

func (q *Queries) CountSessionsBySubjectID(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countSessionsBySubjectID, subjectID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
