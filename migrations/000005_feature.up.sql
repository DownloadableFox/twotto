create table debug_feature (
    guild_id varchar(20) not null,
    name varchar(255) not null,
    enabled boolean not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key (guild_id, name)
);

