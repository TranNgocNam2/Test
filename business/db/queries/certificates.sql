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

-- name: GetCertificateById :one
SELECT *
FROM certificates
WHERE id = sqlc.arg(id);

-- name: GetCertificateByLearnerAndSpecialization :one
SELECT *
FROM certificates
WHERE learner_id = sqlc.arg(learner_id)
AND specialization_id = sqlc.arg(specialization_id)
AND status = sqlc.arg(status)::int;

-- name: GetCertificatesByLearnerAndSubjects :many
SELECT c.*
FROM certificates c
         JOIN specialization_subjects ss
             ON c.specialization_id = ss.specialization_id
                OR c.subject_id = ss.subject_id
WHERE c.learner_id = sqlc.arg(learner_id)
AND c.subject_id = ANY(sqlc.arg(subject_ids)::uuid[])
AND c.status = sqlc.arg(status)::int
GROUP BY c.id, c.subject_id
ORDER BY MAX(ss.index);;

-- name: CreateSubjectCertificate :exec
INSERT INTO certificates (id, learner_id, subject_id, name, status, created_at, class_id)
VALUES ($1, $2, $3, $4, $5, $6, $7);
