ALTER table learner_transcripts
    DROP CONSTRAINT IF EXISTS fk_learner_transcripts_class_learners,
    DROP CONSTRAINT IF EXISTS fk_learner_transcripts_transcript,
    DROP CONSTRAINT IF EXISTS fk_learner_transcripts_updated_by;
ALTER table learner_attendances
    DROP CONSTRAINT IF EXISTS fk_learner_attendances_class_learners,
    DROP CONSTRAINT IF EXISTS fk_learner_attendances_slot,
    DROP CONSTRAINT IF EXISTS unique_learner_attendances;
ALTER table learner_assignments
    DROP CONSTRAINT IF EXISTS fk_learner_assignments_class_learners,
    DROP CONSTRAINT IF EXISTS fk_learner_assignments_assignment,
    DROP CONSTRAINT IF EXISTS unique_learner_assignment;
ALTER table class_learners
    DROP CONSTRAINT IF EXISTS fk_class_learners_learner,
    DROP CONSTRAINT IF EXISTS fk_class_learners_class,
    DROP CONSTRAINT IF EXISTS unique_class_learners;

DROP table IF EXISTS learner_transcripts CASCADE;
DROP table IF EXISTS learner_attendances CASCADE;
DROP table IF EXISTS learner_assignments CASCADE;
DROP table IF EXISTS class_learners CASCADE;
