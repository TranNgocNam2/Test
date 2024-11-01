INSERT INTO subjects(id, code, name, time_per_session, image_link, status, description, created_by)
VALUES('3c1a1849-198a-415e-a136-ad06e114e2bb', 'SEP', 'Lập Trình Machine Learning với Teachable Machine', 2,
'https://notalinkbutalink.com', 1, 'Khung Chương Trình: Lập Trình Machine Learning với Teachable Machine', 'google-oauth2|103166434261305612273');

INSERT INTO subjects(id, code, name, time_per_session, image_link, status, description, created_by)
VALUES('d7fe9772-8ed9-4e96-a274-ebe793d5cb66', 'PRN', 'Programming with .NET', 3,
'https://notalinkbutalink.com', 1,'Loren ipsum', 'google-oauth2|103166434261305612273');

INSERT INTO subjects(id, code, name, time_per_session, image_link, status, description, created_by)
VALUES('010bc611-7b10-4a8b-b0aa-de3342c24641', 'DSA', 'Data structure & Algorithm', 3,
'https://notalinkbutalink.com', 1,'Loren ipsum', 'google-oauth2|103166434261305612273');

INSERT INTO subjects(id, code, name, time_per_session, image_link, status, description, created_by)
VALUES('6838471a-96c9-49c9-b992-281657e04037', 'PRJ', 'Programming with Java', 3,
'https://notalinkbutalink.com', 1,'Loren ipsum', 'google-oauth2|103166434261305612273');


INSERT INTO subjects(id, code, name, time_per_session, image_link, status, description, created_by)
VALUES('aa014905-64d8-456b-af0f-ace63a3ecc2d', 'MAD', 'Discrete Mathematics', 3, 2,
'https://notalinkbutalink.com', 1,'Loren ipsum', 'google-oauth2|103166434261305612273');

INSERT INTO subject_skills(id, subject_id, skill_id)
VALUES  (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '2300b4ec-0a32-4355-bff8-d7a5d6a18bf2'),
        (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '0f3be775-3f5a-4665-bc88-3d72f24b8bf0'),
        (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '3360808a-8b6e-40fa-9cf1-f960a86e5e26');

