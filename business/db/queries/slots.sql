-- name: CreateSlots :copyfrom
INSERT INTO slots (id, session_id, class_id, start_time, end_time, index)
VALUES (sqlc.arg(id), sqlc.arg(session_id), sqlc.arg(class_id), sqlc.arg(start_time),
        sqlc.arg(end_time), sqlc.arg(index));

-- name: GetSlotById :one
SELECT * FROM slots WHERE id = sqlc.arg(id);

-- name: GetSlotByIdAndIndex :one
SELECT * FROM slots WHERE id = sqlc.arg(id) AND index = sqlc.arg(index);

-- name: UpdateSlot :exec
UPDATE slots
SET start_time = sqlc.arg(start_time),
    end_time = sqlc.arg(end_time),
    teacher_id = sqlc.arg(teacher_id)
WHERE id = sqlc.arg(id);

-- name: CheckTeacherTimeOverlap :one
SELECT EXISTS (
    SELECT 1
    FROM slots
    WHERE teacher_id = sqlc.arg(teacher_id)
      AND id <> sqlc.arg(slot_id)
      AND start_time < sqlc.arg(start_time)
      AND end_time > sqlc.arg(end_time)
) AS overlap;

-- name: GetSlotsByClassId :many
SELECT * FROM slots WHERE class_id = sqlc.arg(class_id);

-- name: CountSlotsByClassId :one
SELECT COUNT(*) FROM slots WHERE class_id = sqlc.arg(class_id);

-- name: CountSlotsHaveTeacherByClassId :one
SELECT COUNT(*) FROM slots WHERE class_id = sqlc.arg(class_id) AND teacher_id IS NOT NULL;

-- name: GetSlotByClassIdAndIndex :one
SELECT * FROM slots
    WHERE class_id = sqlc.arg(class_id)
         AND index = sqlc.arg(index);

-- name: UpdateAttendanceCode :exec
UPDATE slots
SET attendance_code = sqlc.arg(attendance_code)
WHERE id = sqlc.arg(id);