CREATE table classes(
    id                      uuid PRIMARY KEY,
    code                    character varying(10) NOT NULL,
    password                character varying(50) NOT NULL,
    name                    character varying(50) NOT NULL,
    subject_id              uuid NOT NULL,
    training_program_id     uuid NOT NULL,
    start_time              timestamp NOT NULL,
    end_time                timestamp NOT NULL,
    created_by              character varying(50) NOT NULL,

    CONSTRAINT fk_class_subject_training_program
        FOREIGN KEY (subject_id, training_program_id)
            REFERENCES training_programs_subjects(subject_id, training_program_id) ON DELETE CASCADE,
    CONSTRAINT fk_class_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_class_subject_training_program UNIQUE (id, subject_id, training_program_id)
);

CREATE table assignments(
    id                  uuid PRIMARY KEY,
    class_id            uuid NOT NULL,
    transcript_id       uuid NOT NULL,

    CONSTRAINT fk_assignment_class
        FOREIGN KEY (class_id)
            REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_transcript
        FOREIGN KEY (transcript_id)
            REFERENCES transcripts(id) ON DELETE CASCADE,

    CONSTRAINT unique_assignment_class_transcript UNIQUE (class_id, transcript_id)
);

CREATE table slots(
    id                  uuid PRIMARY KEY,
    session_id          uuid NOT NULL,
    class_id            uuid NOT NULL,
    start_time          timestamp NOT NULL,
    end_time            timestamp NOT NULL,

    CONSTRAINT fk_slot_sessions
        FOREIGN KEY (session_id)
            REFERENCES sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_slot_class
        FOREIGN KEY (class_id)
            REFERENCES classes(id) ON DELETE CASCADE,

    CONSTRAINT unique_slot_session_class UNIQUE (session_id, class_id)
);