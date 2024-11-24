// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: learner_assignments.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const getAssignmentsByClassLearner = `-- name: GetAssignmentsByClassLearner :many
SELECT id, class_learner_id, assignment_id, grade FROM learner_assignments
    WHERE class_learner_id = $1::uuid
`

func (q *Queries) GetAssignmentsByClassLearner(ctx context.Context, classLearnerID uuid.UUID) ([]LearnerAssignment, error) {
	rows, err := q.db.Query(ctx, getAssignmentsByClassLearner, classLearnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LearnerAssignment
	for rows.Next() {
		var i LearnerAssignment
		if err := rows.Scan(
			&i.ID,
			&i.ClassLearnerID,
			&i.AssignmentID,
			&i.Grade,
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