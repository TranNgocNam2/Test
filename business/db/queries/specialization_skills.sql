-- name: CreateSpecializationSkills :exec
INSERT INTO specialization_skills (specialization_id, skill_id)
SELECT sqlc.arg(specialization_id)::uuid, unnest(sqlc.arg(skill_ids)::uuid[]);

-- name: GetSkillsBySpecialization :many
SELECT skills.id, skills.name
FROM specialization_skills
JOIN skills ON specialization_skills.skill_id = skills.id
WHERE specialization_skills.specialization_id = sqlc.arg(specialization_id)::uuid;

-- name: DeleteSpecializationSkills :exec
DELETE FROM specialization_skills WHERE specialization_id = sqlc.arg(specialization_id)::uuid;