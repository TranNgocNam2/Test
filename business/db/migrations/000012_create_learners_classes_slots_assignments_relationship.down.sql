ALTER table learner_attendances
    DROP CONSTRAINT IF EXISTS fk_learner_attendances_learner,
    DROP CONSTRAINT IF EXISTS fk_learner_attendances_slot,
    DROP CONSTRAINT IF EXISTS unique_learner_attendances;
ALTER table learner_assignments
    DROP CONSTRAINT IF EXISTS fk_learner_assignments_learner,
    DROP CONSTRAINT IF EXISTS fk_learner_assignments_assignment,
    DROP CONSTRAINT IF EXISTS unique_learner_assignment;
ALTER table learners_in_classes
    DROP CONSTRAINT IF EXISTS fk_learners_in_classes_learner,
    DROP CONSTRAINT IF EXISTS fk_learners_in_classes_class,
    DROP CONSTRAINT IF EXISTS unique_learner_in_class;

DROP table IF EXISTS learners_in_classes CASCADE;
DROP table IF EXISTS learner_attendances CASCADE;
DROP table IF EXISTS learner_assignments CASCADE;