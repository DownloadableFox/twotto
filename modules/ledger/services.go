package ledger

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	LedgerManagerKey         = core.NewIdentifier("ledger", "service/manager")
	ErrLedgerManagerNotFound = errors.New("ledger manager not found in context (missing injection)")
)

type LedgerManager interface {
	GetShouldLog(ctx context.Context, guildId string) (bool, error)
	SetShouldLog(ctx context.Context, guildId string, shouldLog bool) error
	GetLogChannel(ctx context.Context, guildId string) (string, error)
	SetLogChannel(ctx context.Context, guildId string, channelId string) error
	LogMessageCreate(ctx context.Context, message *discordgo.Message) error
	LogMessageDelete(ctx context.Context, message *discordgo.Message) error
	LogMessageEdit(ctx context.Context, old *discordgo.Message, new *discordgo.MessageUpdate) error
	LogCustomEvent(ctx context.Context, guildId string, data *discordgo.MessageSend) error
}

type PostgresLedgerManager struct {
	session *discordgo.Session
	pool    *pgxpool.Pool
}

func NewPostgresLedgerManager(session *discordgo.Session, pool *pgxpool.Pool) LedgerManager {
	return &PostgresLedgerManager{pool: pool, session: session}
}

func (m *PostgresLedgerManager) GetShouldLog(ctx context.Context, guildId string) (bool, error) {
	var shouldLog bool
	err := m.pool.QueryRow(ctx, `
		SELECT enabled
		FROM ledger_settings
		WHERE guild_id = $1
	`, guildId).Scan(&shouldLog)
	if err != nil {
		return false, err
	}

	return shouldLog, nil
}

func (m *PostgresLedgerManager) SetShouldLog(ctx context.Context, guildId string, shouldLog bool) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO ledger_settings (guild_id, enabled)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE
		SET enabled = $2
	`, guildId, shouldLog)
	return err
}

func (m *PostgresLedgerManager) GetLogChannel(ctx context.Context, guildId string) (string, error) {
	var channelId string
	err := m.pool.QueryRow(ctx, `
		SELECT log_channel_id
		FROM ledger_settings
		WHERE guild_id = $1
	`, guildId).Scan(&channelId)
	if err != nil {
		return "", err
	}

	return channelId, nil
}

func (m *PostgresLedgerManager) SetLogChannel(ctx context.Context, guildId string, channelId string) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO ledger_settings (guild_id, log_channel_id)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE
		SET log_channel_id = $2
	`, guildId, channelId)
	return err
}

func (m *PostgresLedgerManager) LogCustomEvent(ctx context.Context, guildId string, data *discordgo.MessageSend) error {
	// Get log channel
	channelId, err := m.GetLogChannel(ctx, guildId)
	if err != nil {
		return err
	}

	_, err = m.session.ChannelMessageSendComplex(channelId, data)
	return err
}

func (m *PostgresLedgerManager) LogMessageCreate(ctx context.Context, message *discordgo.Message) error {
	// Logging creates isn't necessary
	return nil
}

func (m *PostgresLedgerManager) LogMessageDelete(ctx context.Context, message *discordgo.Message) error {
	if message.Author.Bot || message.GuildID == "" {
		return nil
	}

	// Check if the guild has logging enabled
	shouldLog, err := m.GetShouldLog(ctx, message.GuildID)
	if err != nil {
		return err
	}

	if !shouldLog {
		return nil
	}

	// Get log channel
	channelId, err := m.GetLogChannel(ctx, message.GuildID)
	if err != nil {
		return err
	}

	if channelId == "" {
		log.Debug().Msgf("[LedgerModule] Guild %s has logging enabled but no log channel set", message.GuildID)
		return nil
	}

	// Send message to log channel
	embed := &discordgo.MessageEmbed{
		Title: "Message Deleted",
		Color: 0x003DFF,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Content",
				Value: message.Content,
			},
			{
				Name:   "Author",
				Value:  message.Author.Mention(),
				Inline: true,
			},
			{
				Name:   "Channel",
				Value:  "<#" + message.ChannelID + ">",
				Inline: true,
			},
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    message.Author.Username + "#" + message.Author.Discriminator,
			IconURL: message.Author.AvatarURL(""),
		},
		Timestamp: message.Timestamp.String(),
	}

	_, err = m.session.ChannelMessageSendEmbed(channelId, embed)
	return err
}

func (m *PostgresLedgerManager) LogMessageEdit(ctx context.Context, old *discordgo.Message, new *discordgo.MessageUpdate) error {
	if new.Author.Bot || new.GuildID == "" {
		return nil
	}

	// Check if the guild has logging enabled
	shouldLog, err := m.GetShouldLog(ctx, new.GuildID)
	if err != nil {
		return err
	}

	if !shouldLog {
		return nil
	}

	// Get log channel
	channelId, err := m.GetLogChannel(ctx, new.GuildID)
	if err != nil {
		return err
	}

	if channelId == "" {
		log.Debug().Msgf("[LedgerModule] Guild %s has logging enabled but no log channel set", new.GuildID)
		return nil
	}

	// Send message to log channel
	embed := &discordgo.MessageEmbed{
		Title: "Message Edited",
		Color: 0xFF003D,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Old Content",
				Value: old.Content,
			},
			{
				Name:  "New Content",
				Value: new.Content,
			},
			{
				Name:   "Author",
				Value:  new.Author.Mention(),
				Inline: true,
			},
			{
				Name:   "Channel",
				Value:  "<#" + new.ChannelID + ">",
				Inline: true,
			},
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    new.Author.Username + "#" + new.Author.Discriminator,
			IconURL: new.Author.AvatarURL(""),
		},
		Timestamp: new.Timestamp.String(),
	}

	_, err = m.session.ChannelMessageSendEmbed(channelId, embed)
	return err
}

/*
func (m *PostgresLedgerManager) LogMessageCreate(ctx context.Context, guildId string, channelId string, messageId string, userId string, content string) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO ledger_messages (message_id, guild_id, channel_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, messageId, guildId, channelId, userId)
	if err != nil {
		return err
	}

	// Log contents
	_, err = m.pool.Exec(ctx, `
		INSERT INTO ledger_contents (message_id, content)
		VALUES ($1, $2)
	`, messageId, content)

	return err
}

func (m *PostgresLedgerManager) LogMessageEdit(ctx context.Context, guildId string, channelId string, messageId string, userId string, content string) error {
	// Check if message exists
	var exists bool
	err := m.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM ledger_messages
			WHERE message_id = $1
		)
	`, messageId).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err := m.pool.Exec(ctx, `
		INSERT INTO ledger_messages (message_id, guild_id, channel_id, user_id)
		VALUES ($1, $2, $3, $4)
		`, messageId, guildId, channelId, userId)
		if err != nil {
			return err
		}
	}

	// Update message to set is_edited to true
	_, err = m.pool.Exec(ctx, `
		UPDATE ledger_messages
		SET is_edited = true
		WHERE message_id = $1
	`, messageId)
	if err != nil {
		return err
	}

	// Log contents
	_, err = m.pool.Exec(ctx, `
		INSERT INTO ledger_contents (message_id, content)
		VALUES ($1, $2)
	`, messageId, content)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresLedgerManager) LogMessageDelete(ctx context.Context, guildId string, channelId string, messageId string) error {
	// Check if message exists
	var exists bool
	err := m.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM ledger_messages
			WHERE message_id = $1
		)
	`, messageId).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err := m.pool.Exec(ctx, `
		INSERT INTO ledger_messages (message_id, guild_id, channel_id)
		VALUES ($1, $2, $3)
	`, messageId, guildId, channelId)
		if err != nil {
			return err
		}
	}

	// Update message to set is_deleted to true
	_, err = m.pool.Exec(ctx, `
		UPDATE ledger_messages
		SET is_deleted = true
		WHERE message_id = $1
	`, messageId)

	return err
}
*/
