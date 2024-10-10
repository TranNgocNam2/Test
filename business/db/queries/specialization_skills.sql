-- name: CreateSpecializationSkills :exec
INSERT INTO specialization_skills (specialization_id, skill_id)
SELECT sqlc.arg(specialization_id)::uuid, unnest(sqlc.arg(skill_ids)::uuid[]);
