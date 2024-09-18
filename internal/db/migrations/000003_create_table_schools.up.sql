CREATE table schools(
    id uuid PRIMARY KEY,
    name character varying(250) NOT NULL,
    address character varying(250) NOT NULL,
    created_by uuid
);