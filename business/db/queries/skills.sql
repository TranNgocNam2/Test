-- name: GetSkillsByIds :many
SELECT * FROM skills WHERE id = ANY(sqlc.arg(skill_ids)::uuid[]);

-- name: GetSkillsBySubjectId :many
SELECT s.id, s.name
FROM skills s
JOIN subject_skills ss ON ss.skill_id = s.id
WHERE ss.subject_id = sqlc.arg(subject_id);

-- name: CreateSkill :exec
INSERT INTO skills (id, name)
VALUES (sqlc.arg(id), sqlc.arg(name));

-- name: UpdateSkill :exec
UPDATE skills
SET name = sqlc.arg(name)
WHERE id = sqlc.arg(id);

-- name: GetSkillById :one
SELECT * FROM skills WHERE id = sqlc.arg(id);

-- name: DeleteSkill :exec
DELETE FROM skills WHERE id = sqlc.arg(id);