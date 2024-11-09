-- name: CreateSpecializationSubjects :exec
INSERT INTO specialization_subjects (specialization_id, subject_id, created_by)
SELECT sqlc.arg(specialization_id)::uuid, unnest(sqlc.arg(subject_ids)::uuid[]), sqlc.arg(created_by)::varchar;

-- name: GetSubjectsBySpecialization :many
SELECT subjects.*
FROM specialization_subjects
JOIN subjects ON specialization_subjects.subject_id = subjects.id
WHERE specialization_subjects.specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: CountSubjectsBySpecializationId :one
SELECT COUNT(*) FROM specialization_subjects
WHERE specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: DeleteSpecializationSubjects :exec
DELETE FROM specialization_subjects WHERE specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: GetSubjectIdsBySpecialization :one
SELECT array_agg(subject_id)::uuid[] as subject_ids
FROM specialization_subjects
WHERE specialization_id = sqlc.arg(specialization_id)::uuid;