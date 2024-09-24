ALTER table learner_specializations
    DROP CONSTRAINT IF EXISTS fk_learner_specialization_learner,
    DROP CONSTRAINT IF EXISTS fk_learner_specialization_specialization;

DROP table IF EXISTS learner_specializations CASCADE;
