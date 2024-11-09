-- name: AddLearnerToSpecialization :exec
INSERT INTO learner_specializations (id, specialization_id, learner_id)
VALUES (uuid_generate_v4(), sqlc.arg(specialization_id), sqlc.arg(learner_id));

-- name: CountLearnerInSpecialization :one
SELECT COUNT(*) FROM learner_specializations
WHERE specialization_id = sqlc.arg(specialization_id)::uuid
AND learner_id = sqlc.arg(learner_id);