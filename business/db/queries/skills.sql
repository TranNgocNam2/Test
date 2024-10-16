-- name: GetSkillsByIDs :many
SELECT * FROM skills WHERE id = ANY(sqlc.arg(skill_ids)::uuid[]);

-- name: GetSkillsBySubjectID :many
SELECT s.id, s.name
FROM skills s
JOIN subject_skills ss ON ss.skill_id = s.id
WHERE ss.subject_id = sqlc.arg(subject_id);
