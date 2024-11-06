// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: class_learners.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const addLearnerToClass = `-- name: AddLearnerToClass :exec
INSERT INTO class_learners (id, class_id, learner_id)
VALUES ($1::uuid, $2::uuid, $3)
`

type AddLearnerToClassParams struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
	LearnerID string    `db:"learner_id" json:"learnerId"`
}

func (q *Queries) AddLearnerToClass(ctx context.Context, arg AddLearnerToClassParams) error {
	_, err := q.db.Exec(ctx, addLearnerToClass, arg.ID, arg.ClassID, arg.LearnerID)
	return err
}

const countLearnersByClassId = `-- name: CountLearnersByClassId :one
SELECT COUNT(*) FROM class_learners WHERE class_id = $1::uuid
`

func (q *Queries) CountLearnersByClassId(ctx context.Context, classID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countLearnersByClassId, classID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
