CREATE table provinces(
    id int PRIMARY KEY,
    name character varying(50) NOT NULL
);

CREATE table districts(
    id int PRIMARY KEY,
    name character varying(50) NOT NULL,
    province_id int NOT NULL
);

CREATE table schools(
    id uuid PRIMARY KEY,
    name character varying(250) NOT NULL,
    address character varying(250) NOT NULL,
    district_id int NOT NULL
);