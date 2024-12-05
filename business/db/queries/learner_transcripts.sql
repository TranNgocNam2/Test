-- name: GenerateLearnerTranscripts :exec
INSERT INTO learner_transcripts (id, class_learner_id, transcript_id)
VALUES (uuid_generate_v4(), sqlc.arg(class_learner_id), UNNEST(sqlc.arg(transcript_ids)::uuid[]));

-- name: GenerateLearnersTranscripts :exec
INSERT INTO learner_transcripts(id, class_learner_id, transcript_id)
SELECT
    uuid_generate_v4(),
    class_learner_id,
    transcript_id
FROM
    UNNEST(sqlc.arg(class_learner_ids)::uuid[]) AS learner_ids(class_learner_id)
        CROSS JOIN
    UNNEST(sqlc.arg(transcript_ids)::uuid[]) AS slot_ids(transcript_id);

-- name: GetLearnerTranscript :one
SELECT * FROM learner_transcripts WHERE class_learner_id = sqlc.arg(class_learner_id) AND transcript_id = sqlc.arg(transcript_id);

-- name: UpdateLearnerTranscriptGrade :exec
UPDATE learner_transcripts
SET grade = sqlc.arg(grade)
WHERE id = sqlc.arg(id)::uuid;
