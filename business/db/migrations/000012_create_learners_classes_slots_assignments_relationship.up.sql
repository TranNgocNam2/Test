CREATE table learners_in_classes(
    id             uuid PRIMARY KEY,
    learner_id     character varying(50) NOT NULL,
    class_id       uuid NOT NULL,

    CONSTRAINT fk_learners_in_classes_learner
        FOREIGN KEY (learner_id)
            REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_learners_in_classes_class
        FOREIGN KEY (class_id)
            REFERENCES classes(id) ON DELETE CASCADE,

    CONSTRAINT unique_learner_in_class UNIQUE (learner_id, class_id)
);

CREATE table learner_attendances(
    id            uuid PRIMARY KEY,
    learner_id    character varying(50) NOT NULL,
    class_id      uuid NOT NULL,
    slot_id       uuid NOT NULL,
    status        int NOT NULL,

    CONSTRAINT fk_learner_attendances_learner
        FOREIGN KEY (learner_id, class_id)
            REFERENCES learners_in_classes(learner_id, class_id) ON DELETE CASCADE,
    CONSTRAINT fk_learner_attendances_slot
        FOREIGN KEY (slot_id)
            REFERENCES slots(id) ON DELETE CASCADE,

    CONSTRAINT unique_learner_attendances UNIQUE (learner_id, slot_id, class_id)
);

CREATE table learner_assignments(
    id              uuid PRIMARY KEY,
    learner_id      character varying(50) NOT NULL,
    class_id        uuid NOT NULL,
    assignment_id   uuid NOT NULL,
    status          int NOT NULL,

    CONSTRAINT fk_learner_assignments_learner
        FOREIGN KEY (learner_id, class_id)
            REFERENCES learners_in_classes(learner_id, class_id) ON DELETE CASCADE,
    CONSTRAINT fk_learner_assignments_assignment
        FOREIGN KEY (assignment_id)
            REFERENCES assignments(id) ON DELETE CASCADE,

    CONSTRAINT unique_learner_assignment UNIQUE (learner_id, assignment_id, class_id)
);