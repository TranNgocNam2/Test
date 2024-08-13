alter table "post"
add column "views" integer not null default 0;

create table "collection"
(
    id         uuid                                not null
        primary key,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    name       text                                not null,
    description text,
    user_id    uuid                                not null
        constraint fk_collection_user
            references "user",
    version    smallint default 1                 not null
);

create table "collection_post"
(
    collection_id uuid not null
        constraint fk_collection_post_collection
            references collection,
    post_id       uuid not null
        constraint fk_collection_post_post
            references post,
    primary key (collection_id, post_id),
    created_at timestamp default CURRENT_TIMESTAMP not null);
