ALTER table schools
    DROP CONSTRAINT IF EXISTS fk_district_schools;
ALTER table districts
    DROP CONSTRAINT IF EXISTS fk_province_districts;
ALTER table verification_learners
    DROP CONSTRAINT IF EXISTS fk_learner_school;

DROP table IF EXISTS provinces CASCADE;
DROP table IF EXISTS districts CASCADE;
DROP table IF EXISTS schools CASCADE;