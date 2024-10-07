ALTER table learners
    DROP CONSTRAINT IF EXISTS fk_learner_user;

ALTER table staffs
    DROP CONSTRAINT IF EXISTS fk_staff_user,
    DROP CONSTRAINT IF EXISTS fk_staff_created_by,
    DROP CONSTRAINT IF EXISTS fk_staff_modified_by;
--     DROP CONSTRAINT IF EXISTS fk_staff_role;

DROP table IF EXISTS learners CASCADE;
DROP table IF EXISTS staffs CASCADE;
DROP table IF EXISTS users CASCADE;

