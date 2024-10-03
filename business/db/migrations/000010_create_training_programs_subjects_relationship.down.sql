ALTER table training_programs_subjects
    DROP CONSTRAINT IF EXISTS fk_training_programs_subjects_training_program,
    DROP CONSTRAINT IF EXISTS fk_training_programs_subjects_subject,
    DROP CONSTRAINT IF EXISTS fk_training_programs_subjects_created_by,
    DROP CONSTRAINT IF EXISTS fk_training_programs_subjects_updated_by,
    DROP CONSTRAINT IF EXISTS unique_training_programs_subjects;

DROP table IF EXISTS training_programs_subjects CASCADE;