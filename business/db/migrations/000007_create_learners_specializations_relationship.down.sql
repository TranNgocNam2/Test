ALTER table learner_specializations
    DROP CONSTRAINT IF EXISTS fk_learners_specializations_learner,
    DROP CONSTRAINT IF EXISTS fk_learners_specializations_specialization,
    DROP CONSTRAINT IF EXISTS unique_learner_specialization;

DROP table IF EXISTS learner_specializations CASCADE;
