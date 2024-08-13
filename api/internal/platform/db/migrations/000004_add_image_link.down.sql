alter table "user"
    drop column avatar_link,
    drop column banner_link;

alter table "artwork"
    drop column image_link,
    drop column processed_image_link;