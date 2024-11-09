-- name: GenerateLearnerAttendance :exec
INSERT INTO learner_attendances(id, class_learner_id, slot_id)
VALUES (uuid_generate_v4(), sqlc.arg(class_learner_id), sqlc.arg(slot_id)::uuid);

-- name: SubmitLearnerAttendance :exec
UPDATE learner_attendances
SET status = sqlc.arg(status)
WHERE id = sqlc.arg(id)::uuid;

-- name: GetLearnerAttendanceByClassLearnerAndSlot :one
SELECT * FROM learner_attendances
    WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid
    AND slot_id = sqlc.arg(slot_id)::uuid;