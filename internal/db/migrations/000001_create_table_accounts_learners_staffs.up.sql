CREATE table accounts (
    id varchar(20) PRIMARY KEY,
    first_name character varying (20) NOT NULL,
    last_name character varying(20) NOT NULL,
    email varchar(20) UNIQUE NOT NULL,
    phone varchar(10) UNIQUE NOT NULL,
    gender int,
    profile_photo varchar(30) NOT NULL,
    status int DEFAULT 1,
    is_deleted bool DEFAULT false
);

CREATE table learners (
    account_id varchar(20) PRIMARY KEY,
    school_id uuid UNIQUE NOT NULL
);

CREATE table staffs (
    account_id varchar(20) PRIMARY KEY,
    role smallint NOT NULL CHECK (role in (1, 2)),
    created_by varchar(20) NOT NULL
);
