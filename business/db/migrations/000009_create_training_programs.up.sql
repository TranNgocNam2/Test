CREATE table training_programs(
    id              uuid PRIMARY KEY,
    name            character varying(100) NOT NULL,
    start_date      date NOT NULL,
    end_date        date NOT NULL,
    created_by      character varying(50) NOT NULL,
    updated_by      character varying(50),
    description     text NOT NULL,
    created_at      timestamp DEFAULT now(),
    updated_at      timestamp,

    CONSTRAINT fk_training_programs_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_training_programs_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE
);