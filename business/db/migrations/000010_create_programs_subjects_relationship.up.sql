CREATE table program_subjects(
    id                      uuid PRIMARY KEY,
    program_id              uuid NOT NULL,
    subject_id              uuid NOT NULL,
    created_by              character varying(50) NOT NULL,
    updated_by              character varying(50),
    created_at              timestamp NOT NULL DEFAULT now(),
    updated_at              timestamp,

    CONSTRAINT fk_program_subjects_program
        FOREIGN KEY (program_id)
            REFERENCES programs(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_subjects_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_subjects_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_subjects_staff_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_program_subjects UNIQUE (program_id, subject_id)
);
