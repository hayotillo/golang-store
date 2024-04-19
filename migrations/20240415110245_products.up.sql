create table if not exists customers
(
    id         char(36) primary key,
    name  varchar(50) not null
);

create trigger set_timestamp
    before update
    on customers
    for each row
execute procedure trigger_set_timestamp();