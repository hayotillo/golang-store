create table if not exists sale_products
(
    id         char(36) primary key,
    sale_id    char(36)    not null,
    product_id char(36)    not null,
    price      numeric(14) not null default 0,
    quantity   numeric(14) not null default 0,

    foreign key (sale_id) references sales (id) on delete cascade,
    foreign key (product_id) references products (id) on delete cascade
);