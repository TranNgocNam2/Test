-- name: InsertSubjectSkill :copyfrom
INSERT INTO subject_skills (id, subject_id, skill_id)
VALUES($1, $2, $3);

