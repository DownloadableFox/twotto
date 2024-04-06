package ledger

import "time"

type LedgerSettings struct {
	GuildId      string    `json:"guild_id"`
	Enabled      bool      `json:"enabled"`
	LogChannelId string    `json:"log_channel_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LedgerMessage struct {
	MessageId string    `json:"message_id"`
	GuildId   string    `json:"guild_id"`
	ChannelId string    `json:"channel_id"`
	UserId    string    `json:"user_id"`
	IsDeleted bool      `json:"is_deleted"`
	IsEdited  bool      `json:"is_edited"`
	CreatedAt time.Time `json:"created_at"`
}

type LedgerContent struct {
	Id        int       `json:"id"`
	MessageId string    `json:"message_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
