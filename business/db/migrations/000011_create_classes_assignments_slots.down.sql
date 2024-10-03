ALTER table classes
    DROP CONSTRAINT IF EXISTS fk_class_subject_training_program,
    DROP CONSTRAINT IF EXISTS fk_class_created_by,
    DROP CONSTRAINT IF EXISTS unique_class_subject_training_program;
ALTER table assignments
    DROP CONSTRAINT IF EXISTS fk_assignment_class,
    DROP CONSTRAINT IF EXISTS fk_assignment_transcript,
    DROP CONSTRAINT IF EXISTS unique_assignment_class_transcript;
ALTER table slots
    DROP CONSTRAINT IF EXISTS fk_slot_sessions,
    DROP CONSTRAINT IF EXISTS fk_slot_class,
    DROP CONSTRAINT IF EXISTS unique_slot_session_class;

DROP table IF EXISTS classes CASCADE;
DROP table IF EXISTS assignments CASCADE;
DROP table IF EXISTS slots CASCADE;