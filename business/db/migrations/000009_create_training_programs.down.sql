ALTER table training_programs
    DROP CONSTRAINT IF EXISTS fk_training_programs_updated_by,
    DROP CONSTRAINT IF EXISTS fk_training_programs_created_by;

DROP table IF EXISTS training_programs CASCADE;