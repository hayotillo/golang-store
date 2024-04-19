create table if not exists visits
(
    id char(36) primary key,
    user_id char(36) not null,
    customer_id char(36) not null,
    status varchar(10) not null default 'wait',
    diagnosis text,
    archive    bool      default false,
    created_at       timestamp    default now(),
    updated_at       timestamp    default now(),

    foreign key (user_id) references users (id) on delete cascade,
    foreign key (customer_id) references customers (id) on delete cascade
);

create trigger set_timestamp
    before update
    on visits
    for each row
execute procedure trigger_set_timestamp();