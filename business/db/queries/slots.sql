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

-- name: UpdateSlotTime :exec
UPDATE slots
SET start_time = sqlc.arg(start_time),
    end_time = sqlc.arg(end_time)
WHERE id = sqlc.arg(id);

-- name: CheckTeacherTimeOverlap :one
SELECT EXISTS (
    SELECT 1
    FROM slots
    WHERE teacher_id = sqlc.arg(teacher_id)
      AND id <> sqlc.arg(slot_id)
      AND start_time IS NOT NULL
      AND end_time IS NOT NULL
      AND NOT (end_time <= sqlc.arg(start_time)
                   OR start_time >= sqlc.arg(end_time))
) AS overlap;

-- name: CheckTeacherTimeOverlapExcludeClass :one
SELECT EXISTS (
    SELECT 1
    FROM slots
    WHERE teacher_id = sqlc.arg(teacher_id)
      AND id <> sqlc.arg(slot_id)
      AND class_id <> slots.class_id
      AND start_time IS NOT NULL
      AND end_time IS NOT NULL
      AND NOT (end_time <= sqlc.arg(start_time)
        OR start_time >= sqlc.arg(end_time))
) AS overlap;

-- name: GetSlotsByClassId :many
SELECT * FROM slots WHERE class_id = sqlc.arg(class_id) ORDER BY index;

-- name: GetConflictingSlotIndexes :one
SELECT STRING_AGG(index+1::TEXT, ',') AS indexes
FROM slots
WHERE class_id = sqlc.arg(class_id)
  AND id <> sqlc.arg(slot_id)
  AND (
       sqlc.arg(new_start_time), sqlc.arg(new_end_time)
          ) OVERLAPS (start_time, end_time);

-- name: CheckSlotOrder :one
SELECT EXISTS (
    SELECT 1
    FROM slots AS s1
             JOIN slots AS s2
                  ON s1.class_id = s2.class_id
                      AND s1.index < s2.index
    WHERE s2.id = sqlc.arg(slot_id)
      AND s2.start_time < s1.end_time
      AND s2.start_time = sqlc.arg(start_time)
      AND s2.end_time = sqlc.arg(end_time)
) AS invalid;

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

-- name: CountCompletedSlotsByClassId :one
SELECT COUNT(*) FROM slots
WHERE class_id = sqlc.arg(class_id)
    AND end_time < now();

-- name: UpdateRecordLink :exec
UPDATE slots
SET record_link = sqlc.arg(record_link)
WHERE id = sqlc.arg(id);
