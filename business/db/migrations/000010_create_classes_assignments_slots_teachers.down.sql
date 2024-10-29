ALTER table classes
    DROP CONSTRAINT IF EXISTS fk_class_staffs_created_by;
ALTER table class_teachers
    DROP CONSTRAINT IF EXISTS fk_class_teachers_teacher,
    DROP CONSTRAINT IF EXISTS fk_class_teachers_class,
    DROP CONSTRAINT IF EXISTS fk_class_teachers_staff_created_by,
    DROP CONSTRAINT IF EXISTS unique_class_teachers;
ALTER table assignments
    DROP CONSTRAINT IF EXISTS fk_assignment_class_teacher,
    DROP CONSTRAINT IF EXISTS fk_assignment_staff_updated_by,
    DROP CONSTRAINT IF EXISTS fk_assignment_transcript,
    DROP CONSTRAINT IF EXISTS unique_assignment_class_teacher_transcript;
ALTER table slots
    DROP CONSTRAINT IF EXISTS fk_slot_sessions,
    DROP CONSTRAINT IF EXISTS fk_slot_class,
    DROP CONSTRAINT IF EXISTS unique_slot_session_class;

DROP table IF EXISTS class_teachers CASCADE;
DROP table IF EXISTS classes CASCADE;
DROP table IF EXISTS assignments CASCADE;
DROP table IF EXISTS slots CASCADE;