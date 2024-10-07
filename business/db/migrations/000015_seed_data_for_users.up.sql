INSERT INTO users (id, full_name, email, phone, gender, auth_role, profile_photo)
VALUES ('google-oauth2|103166434261305612271', 'Trần Ngọc Nam', 'tnnam257@gmail.com', '0886784257', 0, 3,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
        ('google-oauth2|103166434261305612273', 'Trần Ngọc Nam', 'ngocnamsieuquay@gmail.com', '0363832466', 0, 1,
         'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no'),
       ('google-oauth2|103166434261305612274', 'Trần Ngọc Nam', 'ngocnamsieuquay01@gmail.com', '0363832486', 0, 0,
        'https://lh3.googleusercontent.com/a/ACg8ocKU8ncR7STjlTUdE1zws_5-s2sLj9lE9DAUM2sUFu7ho7ReqIFE=s360-c-no');


INSERT INTO staffs (id, role)
VALUES ('google-oauth2|103166434261305612271', 3);

INSERT INTO staffs (id, role, created_by)
VALUES ('google-oauth2|103166434261305612273', 1, 'google-oauth2|103166434261305612271');

INSERT INTO learners (id, school_id)
VALUES ('google-oauth2|103166434261305612274', '3729211e-d63d-4ebf-abf6-9b209c28c2f6');