-- name: GetTeachersByClassID :many
SELECT t.*
FROM class_teachers ct
JOIN users t ON ct.teacher_id = t.id
WHERE ct.class_id = sqlc.arg(class_id)::uuid;