INSERT INTO classes (id, code, password, name, link, subject_id, program_id, start_date, end_date, created_by)
VALUES ('4ce893c8-434d-4592-b307-3b49310c7539', 'FSTEM-01', '12345678',
        'FSTEM-01', 'https://meet.google.com/lookup/1', '3c1a1849-198a-415e-a136-ad06e114e2bb',
        '567c04e1-d067-49f4-867e-c20b33f40991', '2024-10-25', '2024-12-31', 'google-oauth2|103166434261305612280');

INSERT INTO slots (id, session_id, class_id, start_time, end_time, index, attendance_code)
VALUES ('59464179-f2a2-4043-a6dd-062f2b0fff09', '89ac307e-c773-4727-902c-e3b22384e410', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 0, 'QWE123'),
       ('be9e23c3-ee33-4ab9-92e6-a60bbab6abfd', '6d79bb4e-9f6d-4bbb-9bff-191538b7a143', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-11 19:00:00', '2024-11-11 21:00:00', 1, 'QWE123'),
        ('51ca0511-0c8e-418c-a45b-c525e162e6ba', '972e70df-9c0a-4299-86fb-42a4d15d3063', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-13 19:00:00', '2024-11-13 21:00:00', 2, 'QWE123'),
        ('66be8a99-40ff-45c3-91b0-e73658e523bb', 'fb2d510a-d1c8-4e8f-8551-7d14688b7216', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-15 19:00:00', '2024-11-15 21:00:00', 3, 'QWE123'),
        ('e8707ec2-1cb1-424b-abd2-54f7b27d8694', '83b3c33c-52cf-4b83-bbc7-90cd4dca06f2', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-17 19:00:00', '2024-11-17 21:00:00', 4, 'QWE123'),
        ('9d398ef9-5c10-4e15-a147-6309713b9bf2', '8842cbaf-831f-4e04-a13f-cd422daad865', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-19 19:00:00', '2024-11-19 21:00:00', 5, 'QWE123');

INSERT INTO class_learners (id, learner_id, class_id)
VALUES
    ('d6ac21ac-c87e-4e86-b932-69fd73169fcd', 'google-oauth2|103166434261305612269', '4ce893c8-434d-4592-b307-3b49310c7539'),
    ('b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', 'google-oauth2|103166434261305612268', '4ce893c8-434d-4592-b307-3b49310c7539'),
    ('c515b90a-2ea1-4576-9b20-cc751c2238f0', 'google-oauth2|103166434261305612267', '4ce893c8-434d-4592-b307-3b49310c7539'),
    ('bce69116-f98d-4e45-8a96-7f4371e10f01', 'google-oauth2|103166434261305612266', '4ce893c8-434d-4592-b307-3b49310c7539');

INSERT INTO learner_attendances (id, class_learner_id, slot_id, status)
VALUES
    -- Learner 1
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', '59464179-f2a2-4043-a6dd-062f2b0fff09', 0),
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', 'be9e23c3-ee33-4ab9-92e6-a60bbab6abfd', 0),
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', '51ca0511-0c8e-418c-a45b-c525e162e6ba', 0),
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', '66be8a99-40ff-45c3-91b0-e73658e523bb', 0),
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', 'e8707ec2-1cb1-424b-abd2-54f7b27d8694', 0),
    (uuid_generate_v4(), 'd6ac21ac-c87e-4e86-b932-69fd73169fcd', '9d398ef9-5c10-4e15-a147-6309713b9bf2', 0),
    -- Learner 2
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', '59464179-f2a2-4043-a6dd-062f2b0fff09', 0),
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', 'be9e23c3-ee33-4ab9-92e6-a60bbab6abfd', 0),
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', '51ca0511-0c8e-418c-a45b-c525e162e6ba', 0),
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', '66be8a99-40ff-45c3-91b0-e73658e523bb', 0),
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', 'e8707ec2-1cb1-424b-abd2-54f7b27d8694', 0),
    (uuid_generate_v4(), 'b4c98ac1-bd5c-4714-a9f5-d09f8c1680eb', '9d398ef9-5c10-4e15-a147-6309713b9bf2', 0),
    -- Learner 3
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', '59464179-f2a2-4043-a6dd-062f2b0fff09', 0),
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', 'be9e23c3-ee33-4ab9-92e6-a60bbab6abfd', 0),
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', '51ca0511-0c8e-418c-a45b-c525e162e6ba', 0),
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', '66be8a99-40ff-45c3-91b0-e73658e523bb', 0),
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', 'e8707ec2-1cb1-424b-abd2-54f7b27d8694', 0),
    (uuid_generate_v4(), 'c515b90a-2ea1-4576-9b20-cc751c2238f0', '9d398ef9-5c10-4e15-a147-6309713b9bf2', 0),
    -- Learner 4
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', '59464179-f2a2-4043-a6dd-062f2b0fff09', 0),
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', 'be9e23c3-ee33-4ab9-92e6-a60bbab6abfd', 0),
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', '51ca0511-0c8e-418c-a45b-c525e162e6ba', 0),
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', '66be8a99-40ff-45c3-91b0-e73658e523bb', 0),
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', 'e8707ec2-1cb1-424b-abd2-54f7b27d8694', 0),
    (uuid_generate_v4(), 'bce69116-f98d-4e45-8a96-7f4371e10f01', '9d398ef9-5c10-4e15-a147-6309713b9bf2', 0);
