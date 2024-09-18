CREATE table accounts (
    id uuid PRIMARY KEY,
    first_name character varying (10) NOT NULL,
    last_name character varying(10) NOT NULL,
    email varchar(20) UNIQUE NOT NULL,
    phone varchar(10) UNIQUE NOT NULL,
    gender int,
    profile_photo varchar(30) NOT NULL,
    status int DEFAULT 1,
    is_deleted bool DEFAULT false
);

CREATE table learners (
    account_id uuid PRIMARY KEY,
    school_id uuid UNIQUE NOT NULL
);

CREATE table staffs (
    account_id uuid PRIMARY KEY,
    role smallint NOT NULL CHECK (role in (1, 2)),
    created_by uuid NOT NULL
);
