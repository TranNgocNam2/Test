CREATE table program_classes(
    id                      uuid PRIMARY KEY,
    program_id              uuid NOT NULL,
    class_id                uuid NOT NULL,
    created_by              character varying(50) NOT NULL,
    updated_by              character varying(50),
    created_at              timestamp NOT NULL DEFAULT now(),
    updated_at              timestamp,

    CONSTRAINT fk_program_classes_program
        FOREIGN KEY (program_id)
            REFERENCES programs(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_classes_class
        FOREIGN KEY (class_id)
            REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_classes_staff_created_by
        FOREIGN KEY (created_by)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_program_classes_staff_updated_by
        FOREIGN KEY (updated_by)
            REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT unique_program_classes UNIQUE (program_id, class_id)
);
