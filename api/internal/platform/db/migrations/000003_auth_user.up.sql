create table user_role (
                           id int primary key unique not null,
                           role_name varchar(20) not null
);


create table auth_user (
                           id varchar(50) primary key unique not nulL,
                           token text,
                           expiry_date timestamp,
                           created_at timestamp default CURRENT_TIMESTAMP not null,
                           role_id int references user_role(id),
                           user_id uuid unique references "user"(id)
);

alter table user_role
alter column id
add generated  BY DEFAULT  AS identity