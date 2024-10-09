CREATE table specializations(
    id              uuid PRIMARY KEY,
    name            character varying(100) NOT NULL,
    code            character varying(10) NOT NULL,
    time_amount     float NOT NULL,
    image_link      character varying(50) NOT NULL,
    is_drafted      bool DEFAULT false,
    description     text NOT NULL,
    created_by      character varying(50) NOT NULL,
    updated_by      character varying(50),
    created_at      timestamp DEFAULT now(),
    updated_at      timestamp,

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
    sessions_per_week       smallint NOT NULL,
    image_link              character varying(50) NOT NULL,
    status                  smallint  CHECK (status in (0, 1)) NOT NULL,
    description             text NOT NULL,
    created_by              character varying(50) NOT NULL,
    updated_by              character varying(50),
    created_at              timestamp DEFAULT now(),
    updated_at              timestamp,

    CONSTRAINT fk_subject_staff_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subject_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE
);

CREATE table specialization_subjects(
    specialization_id       uuid NOT NULL,
    subject_id              uuid NOT NULL,
    created_by              character varying(50) NOT NULL,

    CONSTRAINT fk_specialization_subjects_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_subjects_subject
        FOREIGN KEY (subject_id)
            REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_subjects_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_specialization_subjects UNIQUE (specialization_id, subject_id)
);

CREATE table skills(
    id      uuid PRIMARY KEY,
    name    character varying(50) NOT NULL
);

CREATE table specialization_skills(
    specialization_id   uuid NOT NULL,
    skill_id            uuid NOT NULL,

    CONSTRAINT fk_specialization_skills_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization_skills_skill
        FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,

    CONSTRAINT unique_specialization_skills UNIQUE (specialization_id, skill_id)
);

CREATE table subject_skills(
    subject_id      uuid NOT NULL,
    skill_id        uuid NOT NULL,

    CONSTRAINT fk_subject_skills_subject
        FOREIGN KEY (subject_id)
            REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_subject_skills_skill
        FOREIGN KEY (skill_id)
            REFERENCES skills(id) ON DELETE CASCADE,

    CONSTRAINT unique_subject_skills UNIQUE (subject_id, skill_id)
);


