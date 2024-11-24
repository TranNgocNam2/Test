CREATE table specializations(
    id              uuid PRIMARY KEY,
    name            character varying(100) NOT NULL,
    code            character varying(10) NOT NULL,
    time_amount     float,
    image_link      text,
    status          smallint DEFAULT 0 CHECK (status in (0, 1, 2)) NOT NULL,
    description     text,
    created_by      character varying(50) NOT NULL,
    updated_by      character varying(50),
    created_at      timestamp with time zone DEFAULT now() NOT NULL,
    updated_at      timestamp with time zone,

    CONSTRAINT fk_specialization_staff_updated_by
        FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_staff_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE table subjects(
    id                      uuid PRIMARY KEY,
    code                    character varying(10) NOT NULL,
    name                    character varying(100) NOT NULL,
    time_per_session        smallint NOT NULL,
    min_pass_grade          real,
    min_attendance          real,
    image_link              text,
    status                  smallint CHECK (status in (0, 1, 2)) DEFAULT 0 NOT NULL,
    description             text,
    created_by              character varying(50) NOT NULL,
    updated_by              character varying(50),
    created_at              timestamp with time zone DEFAULT now() NOT NULL,
    updated_at              timestamp with time zone,
    learner_type            smallint DEFAULT 0,

    CONSTRAINT fk_subject_staff_updated_by
        FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subject_staff_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE table specialization_subjects(
    id                      uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    specialization_id       uuid NOT NULL,
    subject_id              uuid NOT NULL,
    index                   smallint NOT NULL,
    created_by              character varying(50) NOT NULL,

    CONSTRAINT fk_specialization_subjects_specialization
        FOREIGN KEY (specialization_id) REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_subjects_subject
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_subjects_staff_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_specialization_subjects UNIQUE (specialization_id, subject_id)
);

CREATE table skills(
    id      uuid PRIMARY KEY,
    name    character varying(50) NOT NULL
);

CREATE table subject_skills(
    id              uuid PRIMARY KEY,
    subject_id      uuid NOT NULL,
    skill_id        uuid NOT NULL,

    CONSTRAINT fk_subject_skills_subject
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_subject_skills_skill
        FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,

    CONSTRAINT unique_subject_skills UNIQUE (subject_id, skill_id)
);


