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
      AND id <> $2
      AND NOT (end_time <= $3
                   OR start_time >= $4)
) AS overlap
`

type CheckTeacherTimeOverlapParams struct {
	TeacherID *string    `db:"teacher_id" json:"teacherId"`
	SlotID    uuid.UUID  `db:"slot_id" json:"slotId"`
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
}

func (q *Queries) CheckTeacherTimeOverlap(ctx context.Context, arg CheckTeacherTimeOverlapParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkTeacherTimeOverlap,
		arg.TeacherID,
		arg.SlotID,
		arg.StartTime,
		arg.EndTime,
	)
	var overlap bool
	err := row.Scan(&overlap)
	return overlap, err
}

const checkTeacherTimeOverlapExcludeClass = `-- name: CheckTeacherTimeOverlapExcludeClass :one
SELECT EXISTS (
    SELECT 1
    FROM slots
    WHERE teacher_id = $1
      AND id <> $2
      AND class_id <> slots.class_id
      AND NOT (end_time <= $3
        OR start_time >= $4)
) AS overlap
`

type CheckTeacherTimeOverlapExcludeClassParams struct {
	TeacherID *string    `db:"teacher_id" json:"teacherId"`
	SlotID    uuid.UUID  `db:"slot_id" json:"slotId"`
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
}

func (q *Queries) CheckTeacherTimeOverlapExcludeClass(ctx context.Context, arg CheckTeacherTimeOverlapExcludeClassParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkTeacherTimeOverlapExcludeClass,
		arg.TeacherID,
		arg.SlotID,
		arg.StartTime,
		arg.EndTime,
	)
	var overlap bool
	err := row.Scan(&overlap)
	return overlap, err
}

const countCompletedSlotsByClassId = `-- name: CountCompletedSlotsByClassId :one
SELECT COUNT(*) FROM slots
WHERE class_id = $1
    AND end_time < now()
`

func (q *Queries) CountCompletedSlotsByClassId(ctx context.Context, classID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countCompletedSlotsByClassId, classID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countSlotsByClassId = `-- name: CountSlotsByClassId :one
SELECT COUNT(*) FROM slots WHERE class_id = $1
`

func (q *Queries) CountSlotsByClassId(ctx context.Context, classID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countSlotsByClassId, classID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countSlotsHaveTeacherByClassId = `-- name: CountSlotsHaveTeacherByClassId :one
SELECT COUNT(*) FROM slots WHERE class_id = $1 AND teacher_id IS NOT NULL
`

func (q *Queries) CountSlotsHaveTeacherByClassId(ctx context.Context, classID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countSlotsHaveTeacherByClassId, classID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreateSlotsParams struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	SessionID uuid.UUID  `db:"session_id" json:"sessionId"`
	ClassID   uuid.UUID  `db:"class_id" json:"classId"`
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
	Index     int32      `db:"index" json:"index"`
}

const getConflictingSlotIndexes = `-- name: GetConflictingSlotIndexes :one
SELECT STRING_AGG(index::TEXT, ',') AS indexes
FROM slots
WHERE class_id = $1
  AND id <> $2
  AND (
       $3, $4
          ) OVERLAPS (start_time, end_time)
`

type GetConflictingSlotIndexesParams struct {
	ClassID      uuid.UUID   `db:"class_id" json:"classId"`
	SlotID       uuid.UUID   `db:"slot_id" json:"slotId"`
	NewStartTime interface{} `db:"new_start_time" json:"newStartTime"`
	NewEndTime   interface{} `db:"new_end_time" json:"newEndTime"`
}

func (q *Queries) GetConflictingSlotIndexes(ctx context.Context, arg GetConflictingSlotIndexesParams) ([]byte, error) {
	row := q.db.QueryRow(ctx, getConflictingSlotIndexes,
		arg.ClassID,
		arg.SlotID,
		arg.NewStartTime,
		arg.NewEndTime,
	)
	var indexes []byte
	err := row.Scan(&indexes)
	return indexes, err
}

const getSlotByClassIdAndIndex = `-- name: GetSlotByClassIdAndIndex :one
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id, attendance_code, record_link FROM slots
    WHERE class_id = $1
         AND index = $2
`

type GetSlotByClassIdAndIndexParams struct {
	ClassID uuid.UUID `db:"class_id" json:"classId"`
	Index   int32     `db:"index" json:"index"`
}

func (q *Queries) GetSlotByClassIdAndIndex(ctx context.Context, arg GetSlotByClassIdAndIndexParams) (Slot, error) {
	row := q.db.QueryRow(ctx, getSlotByClassIdAndIndex, arg.ClassID, arg.Index)
	var i Slot
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ClassID,
		&i.StartTime,
		&i.EndTime,
		&i.Index,
		&i.TeacherID,
		&i.AttendanceCode,
		&i.RecordLink,
	)
	return i, err
}

const getSlotById = `-- name: GetSlotById :one
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id, attendance_code, record_link FROM slots WHERE id = $1
`

func (q *Queries) GetSlotById(ctx context.Context, id uuid.UUID) (Slot, error) {
	row := q.db.QueryRow(ctx, getSlotById, id)
	var i Slot
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ClassID,
		&i.StartTime,
		&i.EndTime,
		&i.Index,
		&i.TeacherID,
		&i.AttendanceCode,
		&i.RecordLink,
	)
	return i, err
}

const getSlotByIdAndIndex = `-- name: GetSlotByIdAndIndex :one
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id, attendance_code, record_link FROM slots WHERE id = $1 AND index = $2
`

type GetSlotByIdAndIndexParams struct {
	ID    uuid.UUID `db:"id" json:"id"`
	Index int32     `db:"index" json:"index"`
}

func (q *Queries) GetSlotByIdAndIndex(ctx context.Context, arg GetSlotByIdAndIndexParams) (Slot, error) {
	row := q.db.QueryRow(ctx, getSlotByIdAndIndex, arg.ID, arg.Index)
	var i Slot
	err := row.Scan(
		&i.ID,
		&i.SessionID,
		&i.ClassID,
		&i.StartTime,
		&i.EndTime,
		&i.Index,
		&i.TeacherID,
		&i.AttendanceCode,
		&i.RecordLink,
	)
	return i, err
}

const getSlotsByClassId = `-- name: GetSlotsByClassId :many
SELECT id, session_id, class_id, start_time, end_time, index, teacher_id, attendance_code, record_link FROM slots WHERE class_id = $1 ORDER BY index
`

func (q *Queries) GetSlotsByClassId(ctx context.Context, classID uuid.UUID) ([]Slot, error) {
	rows, err := q.db.Query(ctx, getSlotsByClassId, classID)
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
			&i.AttendanceCode,
			&i.RecordLink,
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

const updateAttendanceCode = `-- name: UpdateAttendanceCode :exec
UPDATE slots
SET attendance_code = $1
WHERE id = $2
`

type UpdateAttendanceCodeParams struct {
	AttendanceCode *string   `db:"attendance_code" json:"attendanceCode"`
	ID             uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateAttendanceCode(ctx context.Context, arg UpdateAttendanceCodeParams) error {
	_, err := q.db.Exec(ctx, updateAttendanceCode, arg.AttendanceCode, arg.ID)
	return err
}

const updateRecordLink = `-- name: UpdateRecordLink :exec
UPDATE slots
SET record_link = $1
WHERE id = $2
`

type UpdateRecordLinkParams struct {
	RecordLink *string   `db:"record_link" json:"recordLink"`
	ID         uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) UpdateRecordLink(ctx context.Context, arg UpdateRecordLinkParams) error {
	_, err := q.db.Exec(ctx, updateRecordLink, arg.RecordLink, arg.ID)
	return err
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

const updateSlotTime = `-- name: UpdateSlotTime :exec
UPDATE slots
SET start_time = $1,
    end_time = $2
WHERE id = $3
`

type UpdateSlotTimeParams struct {
	StartTime *time.Time `db:"start_time" json:"startTime"`
	EndTime   *time.Time `db:"end_time" json:"endTime"`
	ID        uuid.UUID  `db:"id" json:"id"`
}

func (q *Queries) UpdateSlotTime(ctx context.Context, arg UpdateSlotTimeParams) error {
	_, err := q.db.Exec(ctx, updateSlotTime, arg.StartTime, arg.EndTime, arg.ID)
	return err
}
