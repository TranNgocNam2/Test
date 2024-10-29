CREATE table classes(
    id                      uuid PRIMARY KEY,
    code                    character varying(10) NOT NULL,
    is_draft                boolean NOT NULL DEFAULT true,
    password                character varying(50) NOT NULL,
    name                    character varying(50) NOT NULL,
    link                    character varying(100) NOT NULL,
    start_time              timestamp NOT NULL,
    end_time                timestamp NOT NULL,
    created_by              character varying(50) NOT NULL,
    created_at              timestamp NOT NULL DEFAULT now(),

    CONSTRAINT fk_class_staffs_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE table class_teachers(
    id             uuid PRIMARY KEY,
    teacher_id     character varying(50) NOT NULL,
    class_id       uuid NOT NULL,
    created_at     timestamp NOT NULL DEFAULT now(),
    created_by     character varying(50) NOT NULL,

    CONSTRAINT fk_class_teachers_teacher
        FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_teachers_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_teachers_staff_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_class_teachers UNIQUE (teacher_id, class_id)
);

CREATE table assignments(
    id                  uuid PRIMARY KEY,
    transcript_id       uuid NOT NULL,
    class_teacher_id    uuid NOT NULL,
    created_at          timestamp NOT NULL DEFAULT now(),
    updated_at          timestamp,
    updated_by          character varying(50),

    CONSTRAINT fk_assignment_class_teacher
        FOREIGN KEY (class_teacher_id) REFERENCES class_teachers(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_staff_updated_by
        FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_assignment_transcript
        FOREIGN KEY (transcript_id) REFERENCES transcripts(id) ON DELETE CASCADE,

    CONSTRAINT unique_assignment_class_teacher_transcript UNIQUE (class_teacher_id, transcript_id)
);

CREATE table slots(
    id                  uuid PRIMARY KEY,
    session_id          uuid NOT NULL,
    class_id            uuid NOT NULL,
    start_time          timestamp NOT NULL,
    end_time            timestamp NOT NULL,

    CONSTRAINT fk_slot_sessions
        FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_slot_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,

    CONSTRAINT unique_slot_session_class UNIQUE (session_id, class_id)
);

