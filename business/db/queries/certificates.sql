-- name: GetCertificationsByLearnerAndSubjects :many
SELECT *
FROM certificates
WHERE learner_id = sqlc.arg(learner_id)
AND subject_id = ANY(sqlc.arg(subject_ids)::uuid[])
AND status = sqlc.arg(status)::int;

-- name: CreateSpecializationCertificate :exec
INSERT INTO certificates (id, learner_id, specialization_id, name, status, created_at)
VALUES (uuid_generate_v4(), sqlc.arg(learner_id), sqlc.arg(specialization_id),
        sqlc.arg(name), sqlc.arg(status), now());