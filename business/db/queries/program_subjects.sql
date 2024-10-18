-- name: CreateProgramSubjects :exec
INSERT INTO program_subjects (id, program_id, subject_id, created_by)
SELECT uuid_generate_v4 (), sqlc.arg(program_id)::uuid, unnest(sqlc.arg(subject_ids)::uuid[]),
       sqlc.arg(created_by)::varchar;

-- name: GetSubjectsByProgram :many
SELECT subjects.id, subjects.name, subjects.code, subjects.image_link, subjects.created_at, subjects.updated_at
FROM program_subjects JOIN subjects ON program_subjects.subject_id = subjects.id
WHERE program_subjects.program_id = sqlc.arg(program_id)::uuid;

-- name: DeleteProgramSubjects :exec
DELETE FROM program_subjects WHERE program_id = sqlc.arg(program_id)::uuid;

-- name: CountSubjectsByProgramID :one
SELECT COUNT(*) FROM program_subjects
WHERE program_id = sqlc.arg(program_id)::uuid;