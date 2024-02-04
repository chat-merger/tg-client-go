pragma foreign_keys = ON;

create table messages
(
    id                  integer primary key autoincrement,
    merger_msg_id       text    not null,
    reply_merger_msg_id text,
    chat_id             integer not null,
    msg_id              integer not null,
    sender_id           integer not null,
    sender_first_name   text    not null,
    kind                integer not null,
    has_media           integer not null, -- bool
    unix_sec            integer not null
);
