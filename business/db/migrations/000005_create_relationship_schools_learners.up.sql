ALTER table learners
ADD CONSTRAINT fk_learner_school
FOREIGN KEY (school_id) REFERENCES schools(id)
ON DELETE CASCADE;