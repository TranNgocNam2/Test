create table "user"
(
    id            uuid                       not null
        primary key,
    username      varchar(50)                         not null
        constraint uni_user_username
            unique,
    avatar_path   varchar(255),
    banner_path   varchar(255),
    name          varchar(50)                         not null,
    date_of_birth timestamp,
    is_ban        boolean   default false             not null,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    version       smallint  default 1                 not null
);

create table "post"
(
    id          uuid                                not null
        primary key,
    title       varchar(500)                        not null,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    is_ban      boolean   default false             not null,
    is_deleted  boolean   default false             not null,
    version     smallint  default 1                 not null,
    user_id     uuid                         not null
        constraint fk_user_posts
            references "user"
);

create table "artwork"
(
    id         uuid                   not null
        primary key,
    image_path varchar(255)           not null,
    processed_image_path varchar(255)           not null,
    is_buyable boolean  default false not null,
    type       smallint               not null,
    is_deleted boolean  default false not null,
    version    smallint default 1     not null,
    artist_id  uuid           not null
        constraint fk_user_artworks
            references "user"
            on update cascade on delete cascade,
    post_id    uuid references post (id)
);

create table "user_downloaded_artwork"
(
    user_id    uuid not null
        constraint fk_user_downloaded_artwork_user
            references "user",
    artwork_id uuid        not null
        constraint fk_user_downloaded_artwork_artwork
            references artwork,
    primary key (user_id, artwork_id)
);

create table "tag"
(
    id              uuid                                not null
        primary key,
    tag_name        varchar(50)                         not null
        constraint uni_tag_tag_name
            unique,
    tag_description text,
    created_at      timestamp default CURRENT_TIMESTAMP not null,
    version         smallint  default 1                 not null
);

create table "post_tags"
(
    post_id uuid not null
        constraint fk_post_tags_post
            references post,
    tag_id  uuid not null
        constraint fk_post_tags_tag
            references tag,
    primary key (post_id, tag_id)
);

create table "user_liked_post"
(
    user_id    uuid                         not null
        constraint fk_user_liked_post_user
            references "user",
    post_id    uuid                                not null
        constraint fk_user_liked_post_post
            references post,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    primary key (user_id, post_id)
);

create table "comment"
(
    id                uuid                                not null
        primary key,
    content           text                                not null,
    created_at        timestamp default CURRENT_TIMESTAMP not null,
    version           smallint  default 1                 not null,
    parent_comment_id uuid
        constraint fk_comment_comments
            references comment,
    post_id           uuid                                not null,
    user_id           uuid not null
        constraint fk_user_comments
            references "user"
            on update cascade on delete cascade
);

create table "wallet"
(
    id             uuid               not null
        primary key,
    omni_balance   bigint   default 0 not null,
    stroke_balance bigint   default 0 not null,
    version        smallint default 1 not null,
    user_id        uuid        not null
        constraint fk_user_wallet
            references "user"
            on update cascade on delete cascade
);

create table "artwork_price"
(
    id         uuid                                not null
        primary key,
    price      bigint,
    from_date  timestamp default CURRENT_TIMESTAMP not null,
    to_date    timestamp,
    version    smallint  default 1                 not null,
    artwork_id uuid                                not null
        constraint fk_artwork_artwork_prices
            references artwork
);

create table "order"
(
    id         uuid                                not null
        primary key,
    order_date timestamp default CURRENT_TIMESTAMP not null,
    amount     bigint    default 0                 not null,
    version    smallint  default 1                 not null,
    buyer_id   uuid                         not null
        constraint fk_user_orders_as_buyer
            references "user"
            on update cascade on delete cascade,
    seller_id  uuid                        not null
        constraint fk_user_orders_as_seller
            references "user"
            on update cascade on delete cascade,
    artwork_id uuid                                not null
        constraint fk_artwork_orders
            references artwork
);

create table "transaction"
(
    id               uuid                                  not null
        primary key,
    type             smallint    default 1                 not null,
    status           smallint    default 1                 not null,
    transaction_date timestamp   default CURRENT_TIMESTAMP not null,
    amount_omni      smallint                              not null,
    amount_stroke    smallint                              not null,
    version          smallint    default 1                 not null,
    buyer_id         uuid                           not null
        constraint fk_user_transactions_as_buyer
            references "user"
            on update cascade on delete cascade,
    seller_id        uuid default NULL:: uuid
        constraint fk_user_transactions_as_seller
            references "user"
            on update cascade on delete cascade,
    order_id         uuid
        constraint fk_order_transactions
            references "order"
            on update cascade on delete cascade
);
