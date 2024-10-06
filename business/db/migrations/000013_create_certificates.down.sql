ALTER table certificates
    DROP CONSTRAINT IF EXISTS fk_certificate_learner,
    DROP CONSTRAINT IF EXISTS fk_certificates_specialization,
    DROP CONSTRAINT IF EXISTS fk_certificates_subject;
--     DROP CONSTRAINT IF EXISTS check_type_specialization_subject;

DROP table IF EXISTS certificates CASCADE;