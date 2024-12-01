CREATE table classes(
    id                      uuid PRIMARY KEY,
    code                    character varying(10) NOT NULL UNIQUE,
    subject_id              uuid NOT NULL,
    program_id              uuid NOT NULL,
    password                character varying(10) NOT NULL,
    name                    character varying(50) NOT NULL,
    link                    text,
    start_date              timestamp with time zone,
    end_date                timestamp with time zone,
    status                  smallint NOT NULL DEFAULT 0,
    type                    smallint NOT NULL DEFAULT 1,
    created_by              character varying(50) NOT NULL,
    created_at              timestamp with time zone NOT NULL DEFAULT now(),
    updated_at              timestamp with time zone,
    updated_by              character varying(50),

    CONSTRAINT fk_class_staffs_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_subject
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_program
        FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_staffs_updated_by
        FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_classes_subject UNIQUE (id, subject_id)
);

CREATE table assignments(
    id                  uuid PRIMARY KEY,
    class_id            uuid NOT NULL,
    question            json NOT NULL,
    deadline            timestamp with time zone,
    status              smallint DEFAULT 0 CHECK (status in (0, 1, 2)) NOT NULL,
    can_overdue         bool DEFAULT false,
    type                smallint DEFAULT 0 CHECK (type in (0, 1, 2)) NOT NULL,

    CONSTRAINT fk_assignment_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);

CREATE table slots(
    id                  uuid PRIMARY KEY,
    session_id          uuid NOT NULL,
    class_id            uuid NOT NULL,
    start_time          timestamp with time zone,
    end_time            timestamp with time zone,
    index               int NOT NULL,
    teacher_id          character varying(50),
    attendance_code     character varying(6) DEFAULT NULL,
    record_link         text,

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
