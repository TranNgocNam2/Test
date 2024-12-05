// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: certificates.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createSpecializationCertificate = `-- name: CreateSpecializationCertificate :exec
INSERT INTO certificates (id, learner_id, specialization_id, name, status, created_at)
VALUES (uuid_generate_v4(), $1, $2,
        $3, $4, now())
`

type CreateSpecializationCertificateParams struct {
	LearnerID        string     `db:"learner_id" json:"learnerId"`
	SpecializationID *uuid.UUID `db:"specialization_id" json:"specializationId"`
	Name             string     `db:"name" json:"name"`
	Status           int32      `db:"status" json:"status"`
}

func (q *Queries) CreateSpecializationCertificate(ctx context.Context, arg CreateSpecializationCertificateParams) error {
	_, err := q.db.Exec(ctx, createSpecializationCertificate,
		arg.LearnerID,
		arg.SpecializationID,
		arg.Name,
		arg.Status,
	)
	return err
}

const getCertificateById = `-- name: GetCertificateById :one
SELECT id, learner_id, specialization_id, subject_id, class_id, name, status, created_at, updated_at, updated_by
FROM certificates
WHERE id = $1
`

func (q *Queries) GetCertificateById(ctx context.Context, id uuid.UUID) (Certificate, error) {
	row := q.db.QueryRow(ctx, getCertificateById, id)
	var i Certificate
	err := row.Scan(
		&i.ID,
		&i.LearnerID,
		&i.SpecializationID,
		&i.SubjectID,
		&i.ClassID,
		&i.Name,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const getCertificationsByLearnerAndSubjects = `-- name: GetCertificationsByLearnerAndSubjects :many
SELECT id, learner_id, specialization_id, subject_id, class_id, name, status, created_at, updated_at, updated_by
FROM certificates
WHERE learner_id = $1
AND subject_id = ANY($2::uuid[])
AND status = $3::int
`

type GetCertificationsByLearnerAndSubjectsParams struct {
	LearnerID  string      `db:"learner_id" json:"learnerId"`
	SubjectIds []uuid.UUID `db:"subject_ids" json:"subjectIds"`
	Status     int32       `db:"status" json:"status"`
}

func (q *Queries) GetCertificationsByLearnerAndSubjects(ctx context.Context, arg GetCertificationsByLearnerAndSubjectsParams) ([]Certificate, error) {
	rows, err := q.db.Query(ctx, getCertificationsByLearnerAndSubjects, arg.LearnerID, arg.SubjectIds, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Certificate
	for rows.Next() {
		var i Certificate
		if err := rows.Scan(
			&i.ID,
			&i.LearnerID,
			&i.SpecializationID,
			&i.SubjectID,
			&i.ClassID,
			&i.Name,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UpdatedBy,
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
