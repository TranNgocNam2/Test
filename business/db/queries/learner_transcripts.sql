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