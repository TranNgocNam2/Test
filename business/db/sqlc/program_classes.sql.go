// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: program_classes.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const countClassesByProgramID = `-- name: CountClassesByProgramID :one
SELECT COUNT(*) FROM program_classes
WHERE program_id = $1::uuid
`

func (q *Queries) CountClassesByProgramID(ctx context.Context, programID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countClassesByProgramID, programID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
