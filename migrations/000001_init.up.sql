-- whitelist module
create table whitelist_settings (
    -- metadata
    guild_id varchar(20) primary key,

    -- settings
    default_role_id varchar(20),
    remove_on_ban boolean not null default false,
    enabled boolean not null default false,

    -- timestamps
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table whitelist_users (
    guild_id varchar(20),
    user_id varchar(20),
    created_at timestamp not null default current_timestamp,
    primary key (guild_id, user_id)
);