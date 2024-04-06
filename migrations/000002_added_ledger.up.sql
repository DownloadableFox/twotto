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

