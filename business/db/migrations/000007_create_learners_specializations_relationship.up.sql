CREATE table learner_specializations(
    learner_id              character varying(50) NOT NULL,
    specialization_id       uuid NOT NULL,
    joined_at               timestamp DEFAULT now(),

    CONSTRAINT fk_learners_specializations_learner
        FOREIGN KEY (learner_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_learners_specializations_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE
);
