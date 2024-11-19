-- Data for learners
INSERT INTO users (id, full_name, email, phone, auth_role, status, profile_photo)
VALUES  ('google-oauth2|103166434261305612270', 'Test Learner 0', 'ngocnamsieuquay01@gmail.com', '0363832486', 0, 1,
        'https://lh3.googleusercontent.com'),
        ('google-oauth2|103166434261305612269', 'Test Learner 1', 'ngocnamsieuquay02@gmail.com', '0363832496', 0, 1,
         'https://lh3.googleusercontent.com'),
        ('google-oauth2|103166434261305612268', 'Test Learner 2', 'ngocnamsieuquay03@gmail.com', '0363832506', 0, 1,
         'https://lh3.googleusercontent.com'),
        ('google-oauth2|103166434261305612267', 'Test Learner 3', 'ngocnamsieuquay04@gmail.com', '0363832516', 0, 1,
         'https://lh3.googleusercontent.com'),
        ('google-oauth2|103166434261305612266', 'Test Learner 4', 'ngocnamsieuquay05@gmail.com', '0363832526', 0, 1,
         'https://lh3.googleusercontent.com');

INSERT INTO users (id, full_name, email, phone, auth_role, profile_photo)
VALUES
    -- Data for admin
    ('google-oauth2|103166434261305612280', 'Test Admin 0', 'tnnam257@gmail.com', '0886784257', 3,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    -- Data for managers
    ('google-oauth2|103166434261305612281', 'Test Manager 0', 'ngocnamsieuquay11@gmail.com', '0886724257', 1,
         'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    -- Data for teachers
    ('google-oauth2|103166434261305612282', 'Test Teacher 0', 'ngocnamsieuquay@gmail.com', '0363832466', 2,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    ('google-oauth2|103166434261305612283', 'Test Teacher 1', 'ngocnamsieuquay12@gmail.com', '0363632456', 2,
         'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no');

-- Data for verification_learners
INSERT INTO verification_learners (id, school_id, learner_id, image_link, status, verified_by, type, verified_at)
VALUES (uuid_generate_v4(), '3729211e-d63d-4ebf-abf6-9b209c28c2f6', 'google-oauth2|103166434261305612270',
        '{"https://lh3.googleusercontent.com"}', 0, 'google-oauth2|103166434261305612282', 1, '2024-11-09 19:00:00'),
       (uuid_generate_v4(), '3729211e-d63d-4ebf-abf6-9b209c28c2f6', 'google-oauth2|103166434261305612269',
        '{"https://lh3.googleusercontent.com"}', 1, 'google-oauth2|103166434261305612282', 0, '2024-11-09 19:00:00'),
       (uuid_generate_v4(), '3729211e-d63d-4ebf-abf6-9b209c28c2f6', 'google-oauth2|103166434261305612268',
        '{"https://lh3.googleusercontent.com"}', 1, 'google-oauth2|103166434261305612282', 0, '2024-11-09 19:00:00'),
       (uuid_generate_v4(), '3729211e-d63d-4ebf-abf6-9b209c28c2f6', 'google-oauth2|103166434261305612267',
        '{"https://lh3.googleusercontent.com"}', 1, 'google-oauth2|103166434261305612282', 0, '2024-11-09 19:00:00'),
       (uuid_generate_v4(), '3729211e-d63d-4ebf-abf6-9b209c28c2f6', 'google-oauth2|103166434261305612266',
        '{"https://lh3.googleusercontent.com"}', 1, 'google-oauth2|103166434261305612282', 1, '2024-11-09 19:00:00');