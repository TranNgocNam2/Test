// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: program_subjects.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countSubjectsByProgramID = `-- name: CountSubjectsByProgramID :one
SELECT COUNT(*) FROM program_subjects
WHERE program_id = $1::uuid
`

func (q *Queries) CountSubjectsByProgramID(ctx context.Context, programID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countSubjectsByProgramID, programID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createProgramSubjects = `-- name: CreateProgramSubjects :exec
INSERT INTO program_subjects (id, program_id, subject_id, created_by)
SELECT uuid_generate_v4 (), $1::uuid, unnest($2::uuid[]),
       $3::varchar
`

type CreateProgramSubjectsParams struct {
	ProgramID  uuid.UUID   `db:"program_id" json:"programId"`
	SubjectIds []uuid.UUID `db:"subject_ids" json:"subjectIds"`
	CreatedBy  string      `db:"created_by" json:"createdBy"`
}

func (q *Queries) CreateProgramSubjects(ctx context.Context, arg CreateProgramSubjectsParams) error {
	_, err := q.db.Exec(ctx, createProgramSubjects, arg.ProgramID, arg.SubjectIds, arg.CreatedBy)
	return err
}

const deleteProgramSubjects = `-- name: DeleteProgramSubjects :exec
DELETE FROM program_subjects WHERE program_id = $1::uuid
`

func (q *Queries) DeleteProgramSubjects(ctx context.Context, programID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteProgramSubjects, programID)
	return err
}

const getSubjectsByProgram = `-- name: GetSubjectsByProgram :many
SELECT subjects.id, subjects.name, subjects.code, subjects.image_link, subjects.created_at, subjects.updated_at
FROM program_subjects JOIN subjects ON program_subjects.subject_id = subjects.id
WHERE program_subjects.program_id = $1::uuid
`

type GetSubjectsByProgramRow struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Code      string     `db:"code" json:"code"`
	ImageLink *string    `db:"image_link" json:"imageLink"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
}

func (q *Queries) GetSubjectsByProgram(ctx context.Context, programID uuid.UUID) ([]GetSubjectsByProgramRow, error) {
	rows, err := q.db.Query(ctx, getSubjectsByProgram, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubjectsByProgramRow
	for rows.Next() {
		var i GetSubjectsByProgramRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Code,
			&i.ImageLink,
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
