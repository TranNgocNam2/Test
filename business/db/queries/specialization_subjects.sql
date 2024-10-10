-- name: CreateSpecializationSubjects :exec
INSERT INTO specialization_subjects (specialization_id, subject_id, created_by)
SELECT sqlc.arg(specialization_id)::uuid, unnest(sqlc.arg(subject_ids)::uuid[]), sqlc.arg(created_by)::varchar;

-- name: GetSubjectsBySpecialization :many
SELECT subjects.id, subjects.name, subjects.code, subjects.image_link, subjects.created_at, subjects.updated_at
FROM specialization_subjects
JOIN subjects ON specialization_subjects.subject_id = subjects.id
WHERE specialization_subjects.specialization_id = sqlc.arg(specialization_id)::uuid;