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

const getClassesByLearnerId = `-- name: GetClassesByLearnerId :many
SELECT id, code, subject_id, program_id, password, name, link, start_date, end_date, status, created_by, created_at, updated_at, updated_by FROM classes
WHERE id IN (SELECT class_id FROM class_learners WHERE learner_id = $1)
`

func (q *Queries) GetClassesByLearnerId(ctx context.Context, learnerID string) ([]Class, error) {
	rows, err := q.db.Query(ctx, getClassesByLearnerId, learnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Class
	for rows.Next() {
		var i Class
		if err := rows.Scan(
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

const getLearnerByClassId = `-- name: GetLearnerByClassId :one
SELECT id, learner_id, class_id FROM class_learners
         WHERE class_id = $1::uuid
           AND learner_id = $2
`

type GetLearnerByClassIdParams struct {
	ClassID   uuid.UUID `db:"class_id" json:"classId"`
	LearnerID string    `db:"learner_id" json:"learnerId"`
}

func (q *Queries) GetLearnerByClassId(ctx context.Context, arg GetLearnerByClassIdParams) (ClassLearner, error) {
	row := q.db.QueryRow(ctx, getLearnerByClassId, arg.ClassID, arg.LearnerID)
	var i ClassLearner
	err := row.Scan(&i.ID, &i.LearnerID, &i.ClassID)
	return i, err
}

const getLearnersByClassId = `-- name: GetLearnersByClassId :many
SELECT u.id, u.full_name, u.email, u.phone, u.auth_role, u.profile_photo, u.status, cl.id AS class_learner_id, s.id AS school_id, s.name AS school_name
FROM users u
        JOIN class_learners cl ON cl.learner_id = u.id
        JOIN classes c ON cl.class_id = c.id
        JOIN verification_learners vl ON u.id = vl.learner_id
        JOIN schools s ON s.id = vl.school_id
WHERE c.id = $1::uuid
`

type GetLearnersByClassIdRow struct {
	ID             string    `db:"id" json:"id"`
	FullName       *string   `db:"full_name" json:"fullName"`
	Email          string    `db:"email" json:"email"`
	Phone          *string   `db:"phone" json:"phone"`
	AuthRole       int16     `db:"auth_role" json:"authRole"`
	ProfilePhoto   *string   `db:"profile_photo" json:"profilePhoto"`
	Status         int32     `db:"status" json:"status"`
	ClassLearnerID uuid.UUID `db:"class_learner_id" json:"classLearnerId"`
	SchoolID       uuid.UUID `db:"school_id" json:"schoolId"`
	SchoolName     string    `db:"school_name" json:"schoolName"`
}

func (q *Queries) GetLearnersByClassId(ctx context.Context, classID uuid.UUID) ([]GetLearnersByClassIdRow, error) {
	rows, err := q.db.Query(ctx, getLearnersByClassId, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLearnersByClassIdRow
	for rows.Next() {
		var i GetLearnersByClassIdRow
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Phone,
			&i.AuthRole,
			&i.ProfilePhoto,
			&i.Status,
			&i.ClassLearnerID,
			&i.SchoolID,
			&i.SchoolName,
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
