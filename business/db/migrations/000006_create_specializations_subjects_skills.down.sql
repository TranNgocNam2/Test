ALTER table specializations
    DROP CONSTRAINT IF EXISTS fk_specialization_staff_created_by,
    DROP CONSTRAINT IF EXISTS fk_specialization_staff_updated_by;
ALTER table specialization_subjects
    DROP CONSTRAINT IF EXISTS fk_specialization_subjects_specialization,
    DROP CONSTRAINT IF EXISTS fk_specialization_subjects_subject,
    DROP CONSTRAINT IF EXISTS fk_specialization_subjects_staff_created_by,
    DROP CONSTRAINT IF EXISTS unique_specialization_subjects;
ALTER table subject_skills
    DROP CONSTRAINT IF EXISTS fk_specializations_skills_specialization,
    DROP CONSTRAINT IF EXISTS fk_specializations_skills_skill,
    DROP CONSTRAINT IF EXISTS unique_specialization_skills;
ALTER table subjects
    DROP CONSTRAINT IF EXISTS fk_subject_staff_created_by,
    DROP CONSTRAINT IF EXISTS fk_subject_staff_updated_by;

DROP table IF EXISTS specializations CASCADE;
DROP table IF EXISTS skills CASCADE;
DROP table IF EXISTS subject_skills CASCADE;
DROP table IF EXISTS specialization_subjects CASCADE;
DROP table IF EXISTS subjects CASCADE;