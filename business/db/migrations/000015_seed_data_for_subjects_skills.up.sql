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

