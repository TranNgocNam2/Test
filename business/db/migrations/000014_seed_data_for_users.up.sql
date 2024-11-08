-- Data for learners
INSERT INTO users (id, full_name, email, phone, gender, auth_role, status, profile_photo,school_id)
VALUES  ('google-oauth2|103166434261305612270', 'Test Learner 0', 'ngocnamsieuquay01@gmail.com', '0363832486', 0, 0, 1,
        'https://lh3.googleusercontent.com', '3729211e-d63d-4ebf-abf6-9b209c28c2f6'),
        ('google-oauth2|103166434261305612269', 'Test Learner 1', 'ngocnamsieuquay02@gmail.com', '0363832496', 0, 0, 1,
         'https://lh3.googleusercontent.com', '3729211e-d63d-4ebf-abf6-9b209c28c2f6'),
        ('google-oauth2|103166434261305612268', 'Test Learner 2', 'ngocnamsieuquay03@gmail.com', '0363832506', 0, 0, 1,
         'https://lh3.googleusercontent.com', '3729211e-d63d-4ebf-abf6-9b209c28c2f6'),
        ('google-oauth2|103166434261305612267', 'Test Learner 3', 'ngocnamsieuquay04@gmail.com', '0363832516', 0, 0, 1,
         'https://lh3.googleusercontent.com', '3729211e-d63d-4ebf-abf6-9b209c28c2f6'),
        ('google-oauth2|103166434261305612266', 'Test Learner 4', 'ngocnamsieuquay05@gmail.com', '0363832526', 0, 0, 1,
         'https://lh3.googleusercontent.com', '3729211e-d63d-4ebf-abf6-9b209c28c2f6');

INSERT INTO users (id, full_name, email, phone, gender, auth_role, profile_photo)
VALUES
    -- Data for admin
    ('google-oauth2|103166434261305612271', 'Test Admin 0', 'tnnam257@gmail.com', '0886784257', 0, 3,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    -- Data for managers
    ('google-oauth2|103166434261305612272', 'Test Manager 0', 'ngocnamsieuquay11@gmail.com', '0886724257', 0, 1,
         'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    -- Data for teachers
    ('google-oauth2|103166434261305612273', 'Test Teacher 0', 'ngocnamsieuquay@gmail.com', '0363832466', 0, 2,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
    ('google-oauth2|103166434261305612274', 'Test Teacher 1', 'ngocnamsieuquay12@gmail.com', '0363632456', 0, 2,
         'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no');