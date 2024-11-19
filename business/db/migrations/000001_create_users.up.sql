CREATE table users (
    id                  character varying(50) PRIMARY KEY,
    full_name           character varying(50),
    email               character varying(50) UNIQUE NOT NULL,
    phone               character varying(10) UNIQUE,
--     gender              smallint  CHECK (gender in (0, 1, 2)),
    auth_role           smallint  NOT NULL DEFAULT 0 CHECK (auth_role in (0, 1, 2, 3)),
    profile_photo       text,
    status              int DEFAULT 1 NOT NULL,
    school_id           uuid,
    image               text [],
    verified_by         character varying(50),

    CONSTRAINT fk_users_verified_by
        FOREIGN KEY (verified_by)
            REFERENCES users(id) ON DELETE CASCADE
);