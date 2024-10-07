ALTER table programs
    DROP CONSTRAINT IF EXISTS fk_programs_staff_updated_by,
    DROP CONSTRAINT IF EXISTS fk_programs_staff_created_by;

DROP table IF EXISTS programs CASCADE;