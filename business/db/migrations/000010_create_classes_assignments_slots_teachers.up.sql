CREATE table classes(
    id                      uuid PRIMARY KEY,
    code                    character varying(10) NOT NULL UNIQUE,
    subject_id              uuid NOT NULL,
    program_id              uuid NOT NULL,
    password                text NOT NULL,
    name                    character varying(50) NOT NULL,
    link                    text,
    start_date              timestamp,
    end_date                timestamp,
    status                  smallint NOT NULL DEFAULT 0,
    created_by              character varying(50) NOT NULL,
    created_at              timestamp NOT NULL DEFAULT now(),

    CONSTRAINT fk_class_staffs_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_subject
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_program
        FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE,

    CONSTRAINT unique_classes_subject UNIQUE (id, subject_id)
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
    start_time          timestamp,
    end_time            timestamp,
    index               int NOT NULL,
    teacher_id          character varying(50),
    attendance_code     character varying(6) DEFAULT NULL,

    CONSTRAINT fk_slot_sessions
        FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    CONSTRAINT fk_slot_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
    CONSTRAINT fk_slot_teacher
        FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_slot_time CHECK (start_time < end_time),

    CONSTRAINT unique_slot_session_class UNIQUE (session_id, class_id)
);

CREATE OR REPLACE FUNCTION update_attendance_status() RETURNS VOID AS $$
BEGIN
    UPDATE learner_attendances
    SET status = 3
    WHERE status = 0
      AND slot_id IN (
        SELECT id
        FROM slots
        WHERE end_time < NOW()
    );
END;
$$ LANGUAGE plpgsql;
