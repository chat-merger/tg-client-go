pragma foreign_keys = ON;

create table store
(
    id    integer primary key autoincrement,
    scope text not null,
    key   text not null,
    value text
);
