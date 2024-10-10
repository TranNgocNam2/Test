-- name: GetSkillsByIDs :many
SELECT * FROM skills WHERE id IN(sqlc.arg(skill_ids)::uuid[]);