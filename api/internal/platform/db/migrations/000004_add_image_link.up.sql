alter table "user"
    add column avatar_link varchar(255),
    add column banner_link varchar(255);

alter table "artwork"
    add column image_link varchar(255),
    add column processed_image_link varchar(255);