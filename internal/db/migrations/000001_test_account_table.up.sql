CREATE TABLE accounts (
    id varchar(36) PRIMARY KEY,
    username varchar(255),
    password varchar(255),
    email varchar(360),
    phone varchar(10),
    address varchar(255)
);

INSERT INTO accounts(id, username, password, email, phone, address)
VALUES ('f6add7e4-42f4-4742-98fb-ab753ceb404a', 'test', '$argon2id$v=19$m=65536,t=1,p=12$fwwuBlycsii+i9H1rhKBQw$S5Mco5TstmnzTZ+hiA2Z3YUdIxWpXrF3W9haQXNChqs', 'test@gmail.com', '0123456789', 'ABC123');
