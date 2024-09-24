CREATE table provinces(
    id int PRIMARY KEY,
    name character varying(50) NOT NULL
);

CREATE table districts(
    id int PRIMARY KEY,
    name character varying(50) NOT NULL,
    province_id int NOT NULL,

    CONSTRAINT fk_province_districts
        FOREIGN KEY (province_id)
            REFERENCES provinces(id) ON DELETE CASCADE
);

CREATE table schools(
    id uuid PRIMARY KEY,
    name character varying(250) NOT NULL,
    address character varying(250) NOT NULL,
    district_id int NOT NULL,
    is_deleted bool DEFAULT false,

    CONSTRAINT fk_district_schools
        FOREIGN KEY (district_id)
            REFERENCES districts(id) ON DELETE CASCADE
);

ALTER table users
    ADD CONSTRAINT fk_learner_school
        FOREIGN KEY (school_id)
            REFERENCES schools(id) ON DELETE CASCADE;