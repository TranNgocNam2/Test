alter table "user"
    add column account_number varchar(20),
    add column bank_name varchar (255),
    add column account_name varchar(255);

create sequence order_code_seq
    start with 29
    increment by 1;

alter table "order"
    add column order_code int not null unique default nextval('order_code_seq');

alter table "order"
add column payment_link varchar(255) default null;

alter table "user"
add column email varchar(255) default null;