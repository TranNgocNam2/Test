alter table "comment"
    drop constraint fk_comment_comments,
    add constraint fk_comment_comments
        foreign key (parent_comment_id)
            references "comment" (id)
            on delete cascade
            on update cascade;