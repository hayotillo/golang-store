create table if not exists product_incoming
(
    id          char(36) primary key,
    user_id     char(36)    not null,
    product_id  char(36)    not null,
    price       numeric(14) not null default 0,
    quantity    numeric(14) not null default 0,
    description text,
    created_at  timestamp            default now(),
    updated_at  timestamp            default now(),

    foreign key (user_id) references users (id) on delete cascade,
    foreign key (product_id) references products (id) on delete cascade
);

create trigger set_timestamp
    before update
    on product_incoming
    for each row
execute procedure trigger_set_timestamp();