INSERT INTO subjects(id, code, name, time_per_session, sessions_per_week, image_link,
                     status, description, created_by, min_pass_grade, min_attendance,
                    total_sessions, learner_type)
VALUES
    ('3c1a1849-198a-415e-a136-ad06e114e2bb', 'ML', 'Lập Trình Machine Learning với Teachable Machine', 2,
     2, 'https://notalinkbutalink.com', 1, 'Lập Trình Machine Learning với Teachable Machine',
     'google-oauth2|103166434261305612280', 5, 5, 6, 1),
    ('d7fe9772-8ed9-4e96-a274-ebe793d5cb66', 'PRN', 'Lập trình .NET', 3,
     3, 'https://example.com/images/programming-dotnet.jpg', 1, 'Khóa học lập trình cơ bản với .NET Framework, bao gồm C# và ASP.NET.',
     'google-oauth2|103166434261305612280', 5, 5, 4, 0),
    ('b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3', 'CS', 'An ninh mạng cơ bản', 3,
     2, 'https://example.com/images/cyber-security.jpg', 1, 'Những khái niệm cơ bản về an ninh mạng, bao gồm mật mã và bảo mật mạng.',
     'google-oauth2|103166434261305612280', 5, 5, 5, 0),
    ('f2a58c91-82d2-4b12-9ed8-54ea6f91bb10', 'WD', 'Phát triển Web với HTML/CSS', 2,
     2, 'https://example.com/images/web-development.jpg', 1, 'Khóa học phát triển web cơ bản với HTML và CSS.',
     'google-oauth2|103166434261305612280', 5, 5, 4, 1),
    ('e3d87aa5-f42d-4b25-9058-4e239f97a13b', 'MLP', 'Nguyên lý Machine Learning', 4,
     3, 'https://example.com/images/machine-learning-principles.jpg', 1, 'Các thuật toán và nguyên lý cơ bản trong Machine Learning.',
     'google-oauth2|103166434261305612280', 5, 5, 5, 0),
    ('12ef93f3-5d28-4e0b-8d43-f1ac9deff501', 'CAL', 'Giải tích I', 3,
     3, 'https://example.com/images/calculus.jpg', 1, 'Các khái niệm cơ bản trong giải tích, giới hạn và đạo hàm.',
     'google-oauth2|103166434261305612280', 5, 5, 5, 1);


INSERT INTO subject_skills(id, subject_id, skill_id)
VALUES  (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '2300b4ec-0a32-4355-bff8-d7a5d6a18bf2'),
        (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '0f3be775-3f5a-4665-bc88-3d72f24b8bf0'),
        (uuid_generate_v4(), '3c1a1849-198a-415e-a136-ad06e114e2bb', '3360808a-8b6e-40fa-9cf1-f960a86e5e26');

INSERT INTO transcripts(id, name, index, min_grade, weight, subject_id)
VALUES
    ('be4d8156-c09f-4a94-9df7-4f84832c1737', 'Progress Test 1', 0, 5, 20, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('544830b0-e12e-43f1-8564-e979c3febcb7', 'Progress Test 2', 1, 5, 20, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('a7599942-d07d-4ddb-adb5-450f0f7e0c54', 'Assignment', 2, 5, 10, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('d712d0a6-6436-4b6d-86c2-ebc2452d5042', 'Final', 3, 5, 50, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('107a51e2-3367-4a7f-aa12-21855f247d11', 'Progress Test 1', 0, 5, 20, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('e55b47d5-5acc-4c17-865e-494b3319a6cb', 'Progress Test 2', 1, 5, 20, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('459bf163-8128-4d9f-bceb-7909f213046e', 'Assignment', 2, 5, 10, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('4a2c1a81-d061-4e73-95f5-aa4d22c832ee', 'Final', 3, 5, 50, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('aac2fb62-b8c4-4418-9cdf-b3f8fde55853', 'Progress Test 1', 0, 5, 20, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('101e2a04-d3fd-4640-9e60-affb95b3ca99', 'Progress Test 2', 1, 5, 20, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('a294d2a7-89ad-489d-92cf-b350ab12025f', 'Assignment', 2, 5, 10, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('20220302-28ad-456f-951c-fc756c3b2490', 'Final', 3, 5, 50, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('428e4991-9e26-4877-9184-339441331a0a', 'Progress Test 1', 0, 5, 20, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('f885b45f-3329-47ba-b2ea-aa983f9740a4', 'Progress Test 2', 1, 5, 20, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('da585c23-2c52-4efa-86cb-0efe11803e01', 'Assignment', 2, 5, 10, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('259da1ae-db26-4b71-978d-f1b6c6140bfa', 'Final', 3, 5, 50, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('87979911-d7fb-4933-b18c-4975c3c3c5ae', 'Progress Test 1', 0, 5, 20, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('b8263ee1-d0c2-4df7-a768-70c84f71e75a', 'Progress Test 2', 1, 5, 20, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('270a81df-d7e7-49cf-bffa-7132647b38e0', 'Assignment', 2, 5, 10, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('6b60794a-ddae-4109-828b-d4326e1728c9', 'Final', 3, 5, 50, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('886636c6-7037-494e-a985-6bfc5fcddd77', 'Progress Test 1', 0, 5, 20, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501'),
    ('23806ade-7be2-4154-b0ad-c686c9e8824c', 'Progress Test 2', 1, 5, 20, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501'),
    ('4fcb831a-4f4a-4a6b-b028-7bd0613c59ff', 'Assignment', 2, 5, 10, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501'),
    ('01504c2f-8a73-4863-afe0-be9fdd0b528c', 'Final', 3, 5, 50, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501');

INSERT INTO sessions(id, name, index, subject_id)
VALUES
    ('54fd2adf-47d6-4680-89a0-31fda67d6818', 'Buoi 1', 0, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('4911d302-128a-458c-8fd1-0d231f0136bb', 'Buoi 2', 1, '3c1a1849-198a-415e-a136-ad06e114e2bb'),
    ('5617b702-0290-4c62-a96b-56a639d0a395', 'Buoi 1', 0, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('520e4631-5418-4498-9984-cb387b8e1e0d', 'Buoi 2', 1, 'd7fe9772-8ed9-4e96-a274-ebe793d5cb66'),
    ('5b6f7396-3062-4b26-afaa-d33264727497', 'Buoi 1', 0, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('f0b4206d-ea9e-4d25-a09d-d3df1f94774d', 'Buoi 2', 1, 'b9f53e13-1e43-4b58-92d8-2f4ad31bbbf3'),
    ('04990dce-f6e8-4bd5-825d-4ddcd6905819', 'Buoi 1', 0, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('5655cb6e-ecaf-45d5-b2e0-19e8a1f49eda', 'Buoi 2', 1, 'f2a58c91-82d2-4b12-9ed8-54ea6f91bb10'),
    ('0c81fd33-736f-4dc5-ba97-677453d15897', 'Buoi 1', 0, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('f99614fa-293c-4f2b-af73-6fc2868dd307', 'Buoi 2', 1, 'e3d87aa5-f42d-4b25-9058-4e239f97a13b'),
    ('5504c448-9e8b-4588-a46a-34de1ee08605', 'Buoi 1', 0, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501'),
    ('5073add0-cb21-445d-99c1-11d7e4242bfc', 'Buoi 2', 1, '12ef93f3-5d28-4e0b-8d43-f1ac9deff501');


INSERT INTO materials(id, name, index, type, data, is_shared, session_id)
VALUES
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '54fd2adf-47d6-4680-89a0-31fda67d6818'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '4911d302-128a-458c-8fd1-0d231f0136bb'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '5617b702-0290-4c62-a96b-56a639d0a395'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '520e4631-5418-4498-9984-cb387b8e1e0d'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '5b6f7396-3062-4b26-afaa-d33264727497'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, 'f0b4206d-ea9e-4d25-a09d-d3df1f94774d'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '04990dce-f6e8-4bd5-825d-4ddcd6905819'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '5655cb6e-ecaf-45d5-b2e0-19e8a1f49eda'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '0c81fd33-736f-4dc5-ba97-677453d15897'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, 'f99614fa-293c-4f2b-af73-6fc2868dd307'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '5504c448-9e8b-4588-a46a-34de1ee08605'),
    (uuid_generate_v4(), 'not_share_1', 0, 'h1', '{"data": "Gioi thieu"}', false, '5073add0-cb21-445d-99c1-11d7e4242bfc');

