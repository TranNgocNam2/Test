-- name: GetTeachersByClassId :many
SELECT t.*
FROM class_teachers ct
JOIN users t ON ct.teacher_id = t.id
WHERE ct.class_id = sqlc.arg(class_id)::uuid AND t.auth_role = 2;

-- name: AddTeacherToClass :exec
INSERT INTO class_teachers (id, teacher_id, class_id, created_by)
SELECT uuid_generate_v4(), unnest(sqlc.arg(teacher_ids)::varchar[]), sqlc.arg(class_id)::uuid, sqlc.arg(created_by);

-- name: RemoveTeacherFromClass :exec
DELETE FROM class_teachers
WHERE class_id = sqlc.arg(class_id)::uuid;