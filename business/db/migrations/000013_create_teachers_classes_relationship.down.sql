ALTER table teachers_in_classes
    DROP CONSTRAINT IF EXISTS fk_teachers_in_classes_teacher,
    DROP CONSTRAINT IF EXISTS fk_teachers_in_classes_class;

DROP table IF EXISTS teachers_in_classes CASCADE;