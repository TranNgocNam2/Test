-- name: CreateSpecializationSubject :exec
INSERT INTO specialization_subjects (specialization_id, subject_id, index, created_by)
VALUES (sqlc.arg(specialization_id)::uuid, sqlc.arg(subject_id)::uuid,
        sqlc.arg(index), sqlc.arg(created_by)::varchar);

-- name: GetSubjectsBySpecialization :many
SELECT subjects.*, specialization_subjects.index
FROM specialization_subjects
JOIN subjects ON specialization_subjects.subject_id = subjects.id
WHERE specialization_subjects.specialization_id = sqlc.arg(specialization_id)::uuid
AND subjects.status = 1 ORDER BY specialization_subjects.index;

-- name: CountSubjectsBySpecializationId :one
SELECT COUNT(*) FROM specialization_subjects
WHERE specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: DeleteSpecializationSubjects :exec
DELETE FROM specialization_subjects WHERE specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: GetSubjectIdsBySpecialization :one
SELECT array_agg(subject_id ORDER BY index)::uuid[] as subject_ids
FROM specialization_subjects
WHERE specialization_id = sqlc.arg(specialization_id)::uuid;