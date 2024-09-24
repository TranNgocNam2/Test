CREATE table specializations(
    id uuid PRIMARY KEY,
    name character varying(100) NOT NULL,
    time_amount int NOT NULL,
    image_link character varying(30) NOT NULL,
    is_draft bool DEFAULT false,
    description text NOT NULL,
    created_by character varying(20) NOT NULL,

    CONSTRAINT fk_specialization_created_by
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE table skills(
    id uuid PRIMARY KEY,
    name character varying(50) NOT NULL
);

CREATE table specializations_skills(
    specialization_id uuid NOT NULL,
    skill_id uuid NOT NULL,

    CONSTRAINT fk_specializations_skills_specialization
        FOREIGN KEY (specialization_id)
            REFERENCES specializations(id) ON DELETE CASCADE,
    CONSTRAINT fk_specializations_skills_skill
        FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);


