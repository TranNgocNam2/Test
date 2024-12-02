CREATE table users (
    id                  character varying(50) PRIMARY KEY,
    full_name           character varying(50),
    email               character varying(50) UNIQUE NOT NULL,
    phone               character varying(10) UNIQUE,
    auth_role           smallint NOT NULL DEFAULT 0 CHECK (auth_role in (0, 1, 2, 3)),
    profile_photo       text,
    status              int DEFAULT 0 NOT NULL,
    is_verified         boolean DEFAULT false NOT NULL,
    school_id           uuid,
    type                smallint
);

CREATE table verification_learners(
    id                  uuid PRIMARY KEY,
    school_id           uuid NOT NULL,
    learner_id          character varying(50) NOT NULL,
    image_link          text[],
    status              smallint DEFAULT 0 NOT NULL,
    verified_by         character varying(50),
    type                smallint DEFAULT 0 NOT NULL,
    verified_at         timestamp without time zone,
    note                text,
    created_at          timestamp without time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_users_learner
        FOREIGN KEY (learner_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_users_verified_by
        FOREIGN KEY (verified_by)
            REFERENCES users(id) ON DELETE CASCADE
);