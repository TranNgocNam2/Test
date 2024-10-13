// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: subjects.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const deleteSubjectSkills = `-- name: DeleteSubjectSkills :exec
DELETE FROM subject_skills WHERE subject_id = $1
`

func (q *Queries) DeleteSubjectSkills(ctx context.Context, subjectID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteSubjectSkills, subjectID)
	return err
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
	ID              uuid.UUID        `db:"id" json:"id"`
	Name            string           `db:"name" json:"name"`
	Code            string           `db:"code" json:"code"`
	Description     string           `db:"description" json:"description"`
	ImageLink       string           `db:"image_link" json:"imageLink"`
	Status          int16            `db:"status" json:"status"`
	TimePerSession  int16            `db:"time_per_session" json:"timePerSession"`
	SessionsPerWeek int16            `db:"sessions_per_week" json:"sessionsPerWeek"`
	CreatedBy       string           `db:"created_by" json:"createdBy"`
	CreatedAt       pgtype.Timestamp `db:"created_at" json:"createdAt"`
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

type InsertSubjectSkillParams struct {
	ID        uuid.UUID `db:"id" json:"id"`
	SubjectID uuid.UUID `db:"subject_id" json:"subjectId"`
	SkillID   uuid.UUID `db:"skill_id" json:"skillId"`
}
