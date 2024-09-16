ALTER table schools
ADD CONSTRAINT fk_staff_schools
FOREIGN KEY (created_by) REFERENCES accounts(id)
ON DELETE CASCADE;

ALTER table learners
ADD CONSTRAINT fk_learner_school
FOREIGN KEY (school_id) REFERENCES schools(id)
ON DELETE CASCADE;