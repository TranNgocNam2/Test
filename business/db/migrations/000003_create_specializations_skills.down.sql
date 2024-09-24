ALTER table specializations
    DROP CONSTRAINT IF EXISTS fk_specialization_created_by;
ALTER table specializations_skills
    DROP CONSTRAINT IF EXISTS fk_specializations_skills_specialization,
    DROP CONSTRAINT IF EXISTS fk_specializations_skills_skill;

DROP table IF EXISTS specializations CASCADE;
DROP table IF EXISTS skills CASCADE;
DROP table IF EXISTS specializations_skills CASCADE;