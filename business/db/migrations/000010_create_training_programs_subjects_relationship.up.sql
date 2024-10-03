CREATE table training_programs_subjects(
    id                      uuid PRIMARY KEY,
    training_program_id     uuid NOT NULL,
    subject_id              uuid NOT NULL,
    created_by              character varying(50) NOT NULL,
    updated_by              character varying(50),
    created_at              timestamp DEFAULT now(),
    updated_at              timestamp,

    CONSTRAINT fk_training_programs_subjects_training_program
        FOREIGN KEY (training_program_id)
            REFERENCES training_programs(id) ON DELETE CASCADE,
    CONSTRAINT fk_training_programs_subjects_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_training_programs_subjects_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_training_programs_subjects_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_training_programs_subjects UNIQUE (training_program_id, subject_id)
);
