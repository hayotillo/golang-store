create table if not exists sales
(
    id char(36) primary key,
    user_id char(36) not null,
    product_id char(36) not null,
    status varchar(10) not null default 'incoming',
    description text,
    archive    bool      default false,
    created_at       timestamp    default now(),
    updated_at       timestamp    default now(),

    foreign key (user_id) references users (id) on delete cascade,
    foreign key (product_id) references products (id) on delete cascade
);

create trigger set_timestamp
    before update
    on sales
    for each row
execute procedure trigger_set_timestamp();