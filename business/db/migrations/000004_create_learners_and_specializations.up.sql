CREATE table learner_specializations(
    learner_id character varying(20) NOT NULL,
    specialization_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),

    CONSTRAINT fk_learners_specializations_learner
        FOREIGN KEY (learner_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_learners_specializations_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE
);
