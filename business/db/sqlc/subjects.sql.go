// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: subjects.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getSubjectsByIDs = `-- name: GetSubjectsByIDs :many
SELECT id, code, name, time_per_session, sessions_per_week, image_link, status, description, created_by, updated_by, created_at, updated_at FROM subjects WHERE id IN($1::uuid[]) AND status = 1
`

func (q *Queries) GetSubjectsByIDs(ctx context.Context, subjectIds []uuid.UUID) ([]Subject, error) {
	rows, err := q.db.Query(ctx, getSubjectsByIDs, subjectIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subject
	for rows.Next() {
		var i Subject
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Name,
			&i.TimePerSession,
			&i.SessionsPerWeek,
			&i.ImageLink,
			&i.Status,
			&i.Description,
			&i.CreatedBy,
			&i.UpdatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const insertSubject = `-- name: InsertSubject :one
INSERT INTO subjects (id, name, code, description, image_link, status,
    time_per_session, sessions_per_week, created_by,
    created_at)
VALUES ($1::uuid, $2, $3, $4,
    $5, $6, $7,
    $8, $9,
    $10)
RETURNING id
`

type InsertSubjectParams struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	Code            string     `db:"code" json:"code"`
	Description     string     `db:"description" json:"description"`
	ImageLink       string     `db:"image_link" json:"imageLink"`
	Status          int16      `db:"status" json:"status"`
	TimePerSession  int16      `db:"time_per_session" json:"timePerSession"`
	SessionsPerWeek int16      `db:"sessions_per_week" json:"sessionsPerWeek"`
	CreatedBy       string     `db:"created_by" json:"createdBy"`
	CreatedAt       *time.Time `db:"created_at" json:"createdAt"`
}

func (q *Queries) InsertSubject(ctx context.Context, arg InsertSubjectParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, insertSubject,
		arg.ID,
		arg.Name,
		arg.Code,
		arg.Description,
		arg.ImageLink,
		arg.Status,
		arg.TimePerSession,
		arg.SessionsPerWeek,
		arg.CreatedBy,
		arg.CreatedAt,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
