ALTER table learners
ADD CONSTRAINT fk_learner_account
FOREIGN KEY (account_id) REFERENCES accounts(id)
ON DELETE CASCADE;

ALTER table staffs

ADD CONSTRAINT fk_staff_account
FOREIGN KEY (account_id) REFERENCES accounts(id)
ON DELETE CASCADE,

ADD CONSTRAINT fk_staff_created_by
FOREIGN KEY (created_by) REFERENCES accounts(id)
ON DELETE CASCADE;


