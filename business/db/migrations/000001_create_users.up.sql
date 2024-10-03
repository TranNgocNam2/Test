CREATE table users (
    id                  character varying(50) PRIMARY KEY,
    full_name           character varying(50) NOT NULL,
    email               character varying(50) UNIQUE NOT NULL,
    phone               character varying(10) UNIQUE NOT NULL,
    gender              smallint NOT NULL DEFAULT 1 CHECK (gender in (1, 2, 3)),
    profile_photo       character varying(30) NOT NULL,
    status              int DEFAULT 1 NOT NULL,
    is_deleted          bool DEFAULT false NOT NULL,
    school_id           uuid DEFAULT NULL,
    role                smallint NOT NULL DEFAULT 1 CHECK (role in (1, 2, 3)),
    created_by          character varying(50),

    CONSTRAINT fk_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE
);