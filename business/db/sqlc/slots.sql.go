// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: slots.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const checkTeacherTimeOverlap = `-- name: CheckTeacherTimeOverlap :one
SELECT EXISTS (
    SELECT 1
    FROM slots
    WHERE teacher_id = $1
      AND start_time < $2
      AND end_time > $3
) AS overlap
`

type CheckTeacherTimeOverlapParams struct {
	TeacherID *string    `db:"teacher_id" json:"teacherId"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
	StartTime *time.Time `db:"start_time" json:"startTime"`
}

func (q *Queries) CheckTeacherTimeOverlap(ctx context.Context, arg CheckTeacherTimeOverlapParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkTeacherTimeOverlap, arg.TeacherID, arg.EndTime, arg.StartTime)
	var overlap bool
	err := row.Scan(&overlap)
	return overlap, err
}

type CreateSlotsParams struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	SessionID uuid.UUID  `db:"session_id" json:"sessionId"`
	ClassID   uuid.UUID  `db:"class_id" json:"classId"`
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
	Index     int32      `db:"index" json:"index"`
}

const getSlotByID = `-- name: GetSlotByID :one
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id FROM slots WHERE id = $1
`

func (q *Queries) GetSlotByID(ctx context.Context, id uuid.UUID) (Slot, error) {
	row := q.db.QueryRow(ctx, getSlotByID, id)
	var i Slot
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ClassID,
		&i.StartTime,
		&i.EndTime,
		&i.Index,
		&i.TeacherID,
	)
	return i, err
}

const getSlotsByClassID = `-- name: GetSlotsByClassID :many
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id FROM slots WHERE class_id = $1
`

func (q *Queries) GetSlotsByClassID(ctx context.Context, classID uuid.UUID) ([]Slot, error) {
	rows, err := q.db.Query(ctx, getSlotsByClassID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Slot
	for rows.Next() {
		var i Slot
		if err := rows.Scan(
			&i.ID,
			&i.SessionID,
			&i.ClassID,
			&i.StartTime,
			&i.EndTime,
			&i.Index,
			&i.TeacherID,
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

const updateSlot = `-- name: UpdateSlot :exec
UPDATE slots
SET start_time = $1,
    end_time = $2,
    teacher_id = $3
WHERE id = $4
`

type UpdateSlotParams struct {
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
	TeacherID *string    `db:"teacher_id" json:"teacherId"`
	ID        uuid.UUID  `db:"id" json:"id"`
}

func (q *Queries) UpdateSlot(ctx context.Context, arg UpdateSlotParams) error {
	_, err := q.db.Exec(ctx, updateSlot,
		arg.StartTime,
		arg.EndTime,
		arg.TeacherID,
		arg.ID,
	)
	return err
}
