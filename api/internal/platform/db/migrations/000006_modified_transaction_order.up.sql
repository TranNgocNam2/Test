alter table "transaction"
    drop column buyer_id,
    drop column seller_id,
    add column user_id uuid default null::uuid constraint fk_user_transactions_as_user
        references "user" on update cascade on delete cascade not null;

alter table "order"
    add column status smallint default 1 not null check (status >= 1 and status <= 3);
