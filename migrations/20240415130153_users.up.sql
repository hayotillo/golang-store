create table if not exists users
(
    id char(36) primary key,
    phone varchar(13) not null unique,
    full_name varchar(30) not null,
    encrypt_password varchar(61),
    token            varchar(255) default '' not null,
    status           varchar(20)  default 'user' not null,
    created_at       timestamp    default now(),
    updated_at       timestamp    default now()
);

create trigger set_timestamp
    before update
    on users
    for each row
execute procedure trigger_set_timestamp();

INSERT INTO public.users (id, phone, full_name, encrypt_password, token, status, created_at, updated_at) VALUES ('5B76F5E7-2149-4937-AA3C-B484808F1376', '+998999068201', 'Mamajanov Xayotillo', '$2a$04$YDB3Z4kliP/UJ4CgdC0Ug.VyynrBRwU3BHYKOpdZJf8r9PjwOCbUu', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVCNzZGNUU3LTIxNDktNDkzNy1BQTNDLUI0ODQ4MDhGMTM3NiJ9.8TGucI6Yo9fv5lrf3CkjpsxjGZDLYc5mV0E8_SISoLw', 'admin', '2023-10-10 11:49:25.820543', '2023-11-02 03:27:45.357111');