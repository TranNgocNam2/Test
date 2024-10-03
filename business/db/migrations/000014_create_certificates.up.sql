CREATE table certificates(
    id                  uuid PRIMARY KEY,
    learner_id          character varying(50) NOT NULL,
    specialization_id   uuid,
    subject_id          uuid,
    name                character varying(50) NOT NULL,
    type                int NOT NULL,
    status              int NOT NULL,
    created_at          timestamp NOT NULL,

    CONSTRAINT fk_certificates_learner
        FOREIGN KEY (learner_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_certificates_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_certificates_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE,

    CONSTRAINT check_type_specialization_subject
        CHECK (
            (type = 1 AND specialization_id IS NOT NULL AND subject_id IS NULL) OR
            (type = 2 AND subject_id IS NOT NULL AND specialization_id IS NULL)
            )
);