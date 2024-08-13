alter table "transaction"
    add column moderator_id uuid default null::uuid constraint fk_user_transactions_as_moderator
    references "user" on update cascade on delete cascade;

create table "user_followed_user" (
    follower_id uuid not null,
    user_id uuid not null,
    followed_at timestamp default CURRENT_TIMESTAMP not null,
    primary key (follower_id, user_id),
    constraint fk_follower
    foreign key (follower_id) references "user" (id) on delete cascade,
    constraint fk_user
    foreign key (user_id) references "user" (id) on delete cascade
);