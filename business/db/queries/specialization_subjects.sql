-- name: CreateSpecializationSubjects :exec
INSERT INTO specialization_subjects (specialization_id, subject_id, created_by)
SELECT sqlc.arg(specialization_id)::uuid, unnest(sqlc.arg(subject_ids)::uuid[]), sqlc.arg(created_by)::varchar;