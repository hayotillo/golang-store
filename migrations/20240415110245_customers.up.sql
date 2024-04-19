create table if not exists customers
(
    id         char(36) primary key,
    full_name  varchar(50) not null,
    birth      date        not null,
    phone      varchar(9),
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create trigger set_timestamp
    before update
    on customers
    for each row
execute procedure trigger_set_timestamp();