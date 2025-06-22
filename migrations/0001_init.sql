create table if not exists users_dialog
(
    id serial primary key,
    user_id bigint unique,
    last_admin_message_id bigint,
    last_user_message_id bigint,
    available bool
);

