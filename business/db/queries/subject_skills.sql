-- name: InsertSubjectSkills :exec
INSERT INTO subject_skills (subject_id, skill_id)
SELECT sqlc.arg(subject_id)::uuid, unnest(sqlc.arg(skill_ids)::uuid[]);
