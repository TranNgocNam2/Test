CREATE table users (
    id character varying(20) PRIMARY KEY,
    full_name character varying(50) NOT NULL,
    email character varying(20) UNIQUE NOT NULL,
    phone character varying(10) UNIQUE NOT NULL,
    gender smallint NOT NULL DEFAULT 1 CHECK (gender in (1, 2, 3)),
    profile_photo character varying(30) NOT NULL,
    status int DEFAULT 1,
    is_deleted bool DEFAULT false,
    school_id uuid UNIQUE NOT NULL,
    role smallint NOT NULL DEFAULT 1 CHECK (role in (1, 2, 3))
--     created_by character varying(20),
--
--     CONSTRAINT fk_staff_created_by
--         FOREIGN KEY (created_by)
--             REFERENCES users(id) ON DELETE CASCADE
);