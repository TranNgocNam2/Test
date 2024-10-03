CREATE table teachers_in_classes(
    id             uuid PRIMARY KEY,
    teacher_id     character varying(50) NOT NULL,
    class_id       uuid NOT NULL,

    CONSTRAINT fk_teachers_in_classes_teacher
        FOREIGN KEY (teacher_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_teachers_in_classes_class
        FOREIGN KEY (class_id)
            REFERENCES classes(id) ON DELETE CASCADE
);