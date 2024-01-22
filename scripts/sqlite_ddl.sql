pragma foreign_keys = ON;

create table message_map
(
    id            integer primary key autoincrement,
    scope         text not null,
    merger_id     text not null,
    controller_id text not null
);
