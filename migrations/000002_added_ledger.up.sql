create table ledger_settings (
    -- metadata
    guild_id varchar(20) primary key,

    -- settings
    enabled boolean not null default false,
    log_channel_id varchar(20),

    -- timestamps
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

/*
create table ledger_messages (
    message_id varchar(20) primary key,
    guild_id varchar(20),
    channel_id varchar(20),
    user_id varchar(20),
    is_deleted boolean not null default false,
    is_edited boolean not null default false,
    created_at timestamp not null default current_timestamp
);

create table ledger_contents (
    id serial primary key,
    message_id varchar(20),
    content text,
    created_at timestamp not null default current_timestamp
    foreign key (message_id) references ledger_messages(message_id)
);
*/