INSERT INTO classes (id, code, password, name, link, subject_id, program_id, start_date, end_date, created_by)
VALUES ('4ce893c8-434d-4592-b307-3b49310c7539', ' FSTEM-01', '$2a$12$OV4gidFEEqo2iEYyrMV3ZOSaqB2oyJpsD4FkSG1ZAexz4z4MQ5fUy',
        'FSTEM-01', 'https://meet.google.com/lookup/1', '3c1a1849-198a-415e-a136-ad06e114e2bb',
        '567c04e1-d067-49f4-867e-c20b33f40991', '2024-10-25', '2024-12-31', 'google-oauth2|103166434261305612272');

INSERT INTO slots (id, session_id, class_id, start_time, end_time, index)
VALUES (uuid_generate_v4(), '89ac307e-c773-4727-902c-e3b22384e410', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 0),
       (uuid_generate_v4(), '6d79bb4e-9f6d-4bbb-9bff-191538b7a143', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 1),
        (uuid_generate_v4(), '972e70df-9c0a-4299-86fb-42a4d15d3063', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 2),
        (uuid_generate_v4(), 'fb2d510a-d1c8-4e8f-8551-7d14688b7216', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 3),
        (uuid_generate_v4(), '83b3c33c-52cf-4b83-bbc7-90cd4dca06f2', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 4),
        (uuid_generate_v4(), '8842cbaf-831f-4e04-a13f-cd422daad865', '4ce893c8-434d-4592-b307-3b49310c7539',
        '2024-11-09 19:00:00', '2024-11-09 21:00:00', 5);