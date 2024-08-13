alter table "wallet"
    drop column omni_balance,
    drop column stroke_balance,
    add column balance bigint default 0 not null;

alter table "transaction"
    drop column amount_omni,
    drop column amount_stroke,
    add column amount bigint default 0 not null;