CREATE table users (
    id                  character varying(50) PRIMARY KEY,
    full_name           character varying(50) NOT NULL,
    email               character varying(50) UNIQUE NOT NULL,
    phone               character varying(10) UNIQUE NOT NULL,
    gender              smallint NOT NULL DEFAULT 0 CHECK (gender in (0, 1, 2)),
    auth_role           smallint  NOT NULL DEFAULT 0 CHECK (auth_role in (0, 1, 2, 3)),
    profile_photo       text NOT NULL,
    status              int DEFAULT 1 NOT NULL
);

CREATE table learners(
    id                  character varying(50) PRIMARY KEY,
    role                smallint NOT NULL DEFAULT 0 CHECK (role in (0)),
    school_id           uuid NOT NULL,

    CONSTRAINT fk_learner_user
        FOREIGN KEY (id)
            REFERENCES users(id) ON DELETE CASCADE
);

CREATE table staffs(
    id                  character varying(50) PRIMARY KEY,
    role                smallint NOT NULL DEFAULT 1 CHECK (role in (1, 2, 3)),
    created_by          character varying(50),
    created_at          timestamp NOT NULL DEFAULT now(),

    CONSTRAINT fk_staff_user
        FOREIGN KEY (id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES staffs(id) ON DELETE CASCADE
);