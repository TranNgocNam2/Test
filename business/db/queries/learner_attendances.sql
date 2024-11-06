-- name: GenerateLearnerAttendance :exec
INSERT INTO learner_attendances(id, class_learner_id, slot_id)
SELECT uuid_generate_v4(), sqlc.arg(class_learner_id), sqlc.arg(slot_id)::uuid;