CREATE table class_learners(
    id             uuid PRIMARY KEY,
    learner_id     character varying(50) NOT NULL,
    class_id       uuid NOT NULL,

    CONSTRAINT fk_class_learners_learner
        FOREIGN KEY (learner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_class_learners_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,

    CONSTRAINT unique_class_learners UNIQUE (learner_id, class_id)
);

CREATE table learner_attendances(
    id                  uuid PRIMARY KEY,
    class_learner_id    uuid NOT NULL,
    slot_id             uuid NOT NULL,
    status              int NOT NULL DEFAULT 0,

    CONSTRAINT fk_learner_attendances_class_learners
        FOREIGN KEY (class_learner_id) REFERENCES class_learners(id) ON DELETE CASCADE,
    CONSTRAINT fk_learner_attendances_slot
        FOREIGN KEY (slot_id) REFERENCES slots(id) ON DELETE CASCADE,

    CONSTRAINT unique_learner_attendances UNIQUE (class_learner_id, slot_id)
);

CREATE table learner_assignments(
    id                  uuid PRIMARY KEY,
    class_learner_id    uuid NOT NULL,
    assignment_id       uuid NOT NULL,
    grade               real NOT NULL,
    data                json,
    grading_status      smallint DEFAULT 0 CHECK (grading_status in (0, 1)) NOT NULL,
    submission_status   smallint DEFAULT 0 CHECK (submission_status in (0, 1, 2)) NOT NULL,

    CONSTRAINT fk_learner_assignments_class_learners
        FOREIGN KEY (class_learner_id)
            REFERENCES class_learners(id) ON DELETE CASCADE,
    CONSTRAINT fk_learner_assignments_assignment
        FOREIGN KEY (assignment_id)
            REFERENCES assignments(id) ON DELETE CASCADE,

    CONSTRAINT unique_learner_assignment UNIQUE (assignment_id, class_learner_id)
);
