ALTER table program_classes
    DROP CONSTRAINT IF EXISTS fk_program_classes_program,
    DROP CONSTRAINT IF EXISTS fk_program_classes_class,
    DROP CONSTRAINT IF EXISTS fk_program_classes_staff_created_by,
    DROP CONSTRAINT IF EXISTS fk_program_classes_staff_updated_by,
    DROP CONSTRAINT IF EXISTS unique_program_classes;

DROP table IF EXISTS program_classes CASCADE;