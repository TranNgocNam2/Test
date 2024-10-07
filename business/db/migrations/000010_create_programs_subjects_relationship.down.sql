ALTER table program_subjects
    DROP CONSTRAINT IF EXISTS fk_program_subjects_program,
    DROP CONSTRAINT IF EXISTS fk_program_subjects_subject,
    DROP CONSTRAINT IF EXISTS fk_program_subjects_staff_created_by,
    DROP CONSTRAINT IF EXISTS fk_program_subjects_staff_updated_by,
    DROP CONSTRAINT IF EXISTS unique_program_subjects;

DROP table IF EXISTS program_subjects CASCADE;