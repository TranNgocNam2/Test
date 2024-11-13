ALTER table classes
    DROP CONSTRAINT IF EXISTS fk_class_staffs_created_by,
    DROP CONSTRAINT IF EXISTS fk_class_subject,
    DROP CONSTRAINT IF EXISTS fk_class_program,
    DROP CONSTRAINT IF EXISTS fk_class_staffs_updated_by,
    DROP CONSTRAINT IF EXISTS unique_classes_subject;
ALTER table assignments
    DROP CONSTRAINT IF EXISTS fk_assignment_teacher,
    DROP CONSTRAINT IF EXISTS fk_assignment_staff_updated_by,
    DROP CONSTRAINT IF EXISTS fk_assignment_transcript,
    DROP CONSTRAINT IF EXISTS fk_assignment_class;
ALTER table slots
    DROP CONSTRAINT IF EXISTS fk_slot_sessions,
    DROP CONSTRAINT IF EXISTS fk_slot_class,
    DROP CONSTRAINT IF EXISTS fk_slot_teacher,
    DROP CONSTRAINT IF EXISTS unique_slot_session_class,
    DROP CONSTRAINT IF EXISTS check_slot_time;

DROP FUNCTION IF EXISTS update_attendance_status();
DROP table IF EXISTS classes CASCADE;
DROP table IF EXISTS assignments CASCADE;
DROP table IF EXISTS slots CASCADE;

