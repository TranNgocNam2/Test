CREATE table sessions(
    id              uuid PRIMARY KEY,
    subject_id      uuid NOT NULL,
    index           int NOT NULL,
    name            character varying(50) NOT NULL,
--     time_amount     int NOT NULL,

    CONSTRAINT fk_session_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE
);

CREATE table materials(
    id              uuid PRIMARY KEY,
    session_id      uuid NOT NULL,
    index           int NOT NULL,
    type            character varying(20) NOT NULL,
    data            json NOT NULL,
    is_shared       bool NOT NULL DEFAULT false,
    name            character varying(100),

    CONSTRAINT fk_material_session
        FOREIGN KEY (session_id)
            REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE table transcripts
(
    id              uuid PRIMARY KEY,
    subject_id      uuid NOT NULL,
    name            character varying(50) NOT NULL,
    index           int NOT NULL,
    min_grade       float NOT NULL,
    weight          float NOT NULL,

    CONSTRAINT fk_transcript_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE
);
