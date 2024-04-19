create table if not exists products
(
    id         char(36) primary key,
    name  varchar(50) not null,
    constraint products_unique unique (name)
);

create trigger set_timestamp
    before update
    on products
    for each row
execute procedure trigger_set_timestamp();