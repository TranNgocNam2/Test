ALTER TABLE verification_learners
    DROP CONSTRAINT IF EXISTS fk_users_learner,
    DROP CONSTRAINT IF EXISTS fk_users_verified_by;

DROP table IF EXISTS verification_learners CASCADE;
DROP table IF EXISTS users CASCADE;