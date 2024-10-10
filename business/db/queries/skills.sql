-- name: GetSkillsByIDs :many
SELECT * FROM skills WHERE id = ANY(sqlc.arg(skill_ids)::uuid[]);