// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: classes.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countClassesByProgramId = `-- name: CountClassesByProgramId :one
SELECT COUNT(*) FROM classes WHERE program_id = $1::uuid
`

func (q *Queries) CountClassesByProgramId(ctx context.Context, programID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countClassesByProgramId, programID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createClass = `-- name: CreateClass :exec
INSERT INTO classes (id, code, password, name, subject_id, program_id, link, start_date, end_date, created_by)
VALUES ($1::uuid, $2, $3,
        $4, $5::uuid, $6::uuid,
        $7, $8, $9, $10)
`

type CreateClassParams struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Code      string     `db:"code" json:"code"`
	Password  string     `db:"password" json:"password"`
	Name      string     `db:"name" json:"name"`
	SubjectID uuid.UUID  `db:"subject_id" json:"subjectId"`
	ProgramID uuid.UUID  `db:"program_id" json:"programId"`
	Link      *string    `db:"link" json:"link"`
	StartDate *time.Time `db:"start_date" json:"startDate"`
	EndDate   *time.Time `db:"end_date" json:"endDate"`
	CreatedBy string     `db:"created_by" json:"createdBy"`
}

func (q *Queries) CreateClass(ctx context.Context, arg CreateClassParams) error {
	_, err := q.db.Exec(ctx, createClass,
		arg.ID,
		arg.Code,
		arg.Password,
		arg.Name,
		arg.SubjectID,
		arg.ProgramID,
		arg.Link,
		arg.StartDate,
		arg.EndDate,
		arg.CreatedBy,
	)
	return err
}

const deleteClass = `-- name: DeleteClass :exec
DELETE FROM classes WHERE id = $1::uuid
`

func (q *Queries) DeleteClass(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteClass, id)
	return err
}

const getClassById = `-- name: GetClassById :one
SELECT id, code, subject_id, program_id, password, name, link, start_date, end_date, status, created_by, created_at FROM classes WHERE id = $1::uuid
`

func (q *Queries) GetClassById(ctx context.Context, id uuid.UUID) (Class, error) {
	row := q.db.QueryRow(ctx, getClassById, id)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.SubjectID,
		&i.ProgramID,
		&i.Password,
		&i.Name,
		&i.Link,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const getClassCompletedByCode = `-- name: GetClassCompletedByCode :one
SELECT id, code, subject_id, program_id, password, name, link, start_date, end_date, status, created_by, created_at FROM classes WHERE code = $1 AND status = 1
`

func (q *Queries) GetClassCompletedByCode(ctx context.Context, code string) (Class, error) {
	row := q.db.QueryRow(ctx, getClassCompletedByCode, code)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.SubjectID,
		&i.ProgramID,
		&i.Password,
		&i.Name,
		&i.Link,
		&i.StartDate,
		&i.EndDate,
		&i.Status,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const softDeleteClass = `-- name: SoftDeleteClass :exec
UPDATE classes
SET status = 2
WHERE id = $1::uuid
`

func (q *Queries) SoftDeleteClass(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, softDeleteClass, id)
	return err
}

const updateActiveClass = `-- name: UpdateActiveClass :exec
UPDATE classes
SET status = $1,
    start_date = $2,
    end_date = $3
WHERE id = $4::uuid
`

type UpdateActiveClassParams struct {
	Status    int16      `db:"status" json:"status"`
	StartDate *time.Time `db:"start_date" json:"startDate"`
	EndDate   *time.Time `db:"end_date" json:"endDate"`
	ID        uuid.UUID  `db:"id" json:"id"`
}

func (q *Queries) UpdateActiveClass(ctx context.Context, arg UpdateActiveClassParams) error {
	_, err := q.db.Exec(ctx, updateActiveClass,
		arg.Status,
		arg.StartDate,
		arg.EndDate,
		arg.ID,
	)
	return err
}

const updateClass = `-- name: UpdateClass :exec
UPDATE classes
SET name = $1,
    code = $2,
    password = $3
WHERE id = $4::uuid
`

type UpdateClassParams struct {
	Name     string    `db:"name" json:"name"`
	Code     string    `db:"code" json:"code"`
	Password string    `db:"password" json:"password"`
	ID       uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateClass(ctx context.Context, arg UpdateClassParams) error {
	_, err := q.db.Exec(ctx, updateClass,
		arg.Name,
		arg.Code,
		arg.Password,
		arg.ID,
	)
	return err
}
