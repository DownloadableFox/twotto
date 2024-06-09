package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	LedgerManagerKey         = core.NewIdentifier("ledger", "service/manager")
	ErrLedgerManagerNotFound = errors.New("ledger manager not found in context (missing injection)")
)

type LedgerRepository interface {
	GetLedgerSettings(ctx context.Context, guildId string) (*LedgerSettings, error)
	GetAllLedgerSettings(ctx context.Context, limit int, page int) ([]*LedgerSettings, error)
	CreateLedgerSettings(ctx context.Context, settings *LedgerSettings) error
	UpdateLedgerSettings(ctx context.Context, settings *LedgerSettings) error
	DeleteLedgerSettings(ctx context.Context, guildId string) error
	GetMessage(ctx context.Context, messageId string) (*LedgerMessage, error)
	GetMessages(ctx context.Context, guildId string, limit int, page int) ([]*LedgerMessage, error)
	CreateMessage(ctx context.Context, message *LedgerMessage) error
	UpdateMessage(ctx context.Context, message *LedgerMessage) error
	DeleteMessage(ctx context.Context, messageId string) error
	GetMessageContent(ctx context.Context, contentId int) (*LedgerContent, error)
	GetMessageContents(ctx context.Context, messageId string) ([]*LedgerContent, error)
	CreateMessageContent(ctx context.Context, content *LedgerContent) error
	UpdateMessageContent(ctx context.Context, content *LedgerContent) error
	DeleteMessageContent(ctx context.Context, contentId int) error
}

type LedgerPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewLedgerPostgresRepository(pool *pgxpool.Pool) LedgerRepository {
	return &LedgerPostgresRepository{pool: pool}
}

func (r *LedgerPostgresRepository) GetLedgerSettings(ctx context.Context, guildId string) (*LedgerSettings, error) {
	var settings LedgerSettings
	err := r.pool.QueryRow(ctx, `
		SELECT guild_id, enabled, log_channel_id, created_at, updated_at
		FROM ledger_settings
		WHERE guild_id = $1	
	`, guildId).Scan(&settings.GuildId, &settings.Enabled, &settings.LogChannelId, &settings.CreatedAt, &settings.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ledgerSettings := &LedgerSettings{GuildId: guildId, Enabled: false, LogChannelId: ""}
			err := r.CreateLedgerSettings(ctx, ledgerSettings)
			if err != nil {
				return nil, err
			}

			return ledgerSettings, nil
		}

		return nil, err
	}

	return &settings, nil
}

func (r *LedgerPostgresRepository) GetAllLedgerSettings(ctx context.Context, limit int, page int) ([]*LedgerSettings, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT guild_id, enabled, log_channel_id, created_at, updated_at
		FROM ledger_settings
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, page*limit)
	if err != nil {
		return nil, err
	}

	var settings []*LedgerSettings
	for rows.Next() {
		var setting LedgerSettings
		if err := rows.Scan(&setting.GuildId, &setting.Enabled, &setting.LogChannelId, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, &setting)
	}

	return settings, nil
}

func (r *LedgerPostgresRepository) CreateLedgerSettings(ctx context.Context, settings *LedgerSettings) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO ledger_settings (guild_id, enabled, log_channel_id)
		VALUES ($1, $2, $3)
	`, settings.GuildId, settings.Enabled, settings.LogChannelId)
	return err
}

func (r *LedgerPostgresRepository) UpdateLedgerSettings(ctx context.Context, settings *LedgerSettings) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE ledger_settings
		SET enabled = $1, log_channel_id = $2
		WHERE guild_id = $3
	`, settings.Enabled, settings.LogChannelId, settings.GuildId)
	return err
}

func (r *LedgerPostgresRepository) DeleteLedgerSettings(ctx context.Context, guildId string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM ledger_settings
		WHERE guild_id = $1
	`, guildId)
	return err
}

func (r *LedgerPostgresRepository) GetMessage(ctx context.Context, messageId string) (*LedgerMessage, error) {
	var message LedgerMessage
	err := r.pool.QueryRow(ctx, `
		SELECT message_id, guild_id, channel_id, user_id, is_deleted, is_edited, created_at
		FROM ledger_messages
		WHERE message_id = $1
	`, messageId).Scan(&message.MessageId, &message.GuildId, &message.ChannelId, &message.UserId, &message.IsDeleted, &message.IsEdited, &message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *LedgerPostgresRepository) GetMessages(ctx context.Context, guildId string, limit int, page int) ([]*LedgerMessage, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT message_id, guild_id, channel_id, user_id, is_deleted, is_edited, created_at
		FROM ledger_messages
		WHERE guild_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, guildId, limit, page*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*LedgerMessage
	for rows.Next() {
		var message LedgerMessage
		if err := rows.Scan(&message.MessageId, &message.GuildId, &message.ChannelId, &message.UserId, &message.IsDeleted, &message.IsEdited, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *LedgerPostgresRepository) CreateMessage(ctx context.Context, message *LedgerMessage) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO ledger_messages (message_id, guild_id, channel_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, message.MessageId, message.GuildId, message.ChannelId, message.UserId)
	return err
}

func (r *LedgerPostgresRepository) UpdateMessage(ctx context.Context, message *LedgerMessage) error {
	// if message doesn't exist create before updating
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM ledger_messages
			WHERE message_id = $1
		)
	`, message.MessageId).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return r.CreateMessage(ctx, message)
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE ledger_messages
		SET is_deleted = $1, is_edited = $2
		WHERE message_id = $3
	`, message.IsDeleted, message.IsEdited, message.MessageId)
	return err
}

func (r *LedgerPostgresRepository) DeleteMessage(ctx context.Context, messageId string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM ledger_messages
		WHERE message_id = $1
	`, messageId)
	return err
}

func (r *LedgerPostgresRepository) GetMessageContent(ctx context.Context, contentId int) (*LedgerContent, error) {
	var content LedgerContent
	err := r.pool.QueryRow(ctx, `
		SELECT id, message_id, content, created_at
		FROM ledger_contents
		WHERE id = $1
	`, contentId).Scan(&content.Id, &content.MessageId, &content.Content, &content.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func (r *LedgerPostgresRepository) GetMessageContents(ctx context.Context, messageId string) ([]*LedgerContent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, message_id, content, created_at
		FROM ledger_contents
		WHERE message_id = $1
	`, messageId)
	if err != nil {
		return nil, err
	}

	var contents []*LedgerContent
	for rows.Next() {
		var content LedgerContent
		if err := rows.Scan(&content.Id, &content.MessageId, &content.Content, &content.CreatedAt); err != nil {
			return nil, err
		}
		contents = append(contents, &content)
	}

	return contents, nil
}

func (r *LedgerPostgresRepository) CreateMessageContent(ctx context.Context, content *LedgerContent) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO ledger_contents (message_id, content)
		VALUES ($1, $2)
	`, content.MessageId, content.Content)
	return err
}

func (r *LedgerPostgresRepository) UpdateMessageContent(ctx context.Context, content *LedgerContent) error {
	// if content doesn't exist create before updating
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM ledger_contents
			WHERE id = $1
		)
	`, content.Id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return r.CreateMessageContent(ctx, content)
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE ledger_contents
		SET content = $1
		WHERE id = $2
	`, content.Content, content.Id)
	return err
}

func (r *LedgerPostgresRepository) DeleteMessageContent(ctx context.Context, contentId int) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM ledger_contents
		WHERE id = $1
	`, contentId)
	return err
}

type LedgerManager interface {
	GetShouldLog(ctx context.Context, guildId string) (bool, error)
	SetShouldLog(ctx context.Context, guildId string, shouldLog bool) error
	GetLogChannel(ctx context.Context, guildId string) (string, error)
	SetLogChannel(ctx context.Context, guildId string, channelId string) error
	LogMessageCreate(ctx context.Context, message *discordgo.Message) error
	LogMessageDelete(ctx context.Context, message *discordgo.Message) error
	LogMessageEdit(ctx context.Context, message *discordgo.MessageUpdate) error
	LogCustomEvent(ctx context.Context, guildId string, data *discordgo.MessageSend) error
}

type RepoLedgerManager struct {
	repo    LedgerRepository
	session *discordgo.Session
}

func NewRepoLedgerManager(repo LedgerRepository, session *discordgo.Session) LedgerManager {
	return &RepoLedgerManager{
		repo:    repo,
		session: session,
	}
}

func (m *RepoLedgerManager) GetShouldLog(ctx context.Context, guildId string) (bool, error) {
	settings, err := m.repo.GetLedgerSettings(ctx, guildId)
	if err != nil {
		return false, err
	}

	return settings.Enabled, nil
}

func (m *RepoLedgerManager) SetShouldLog(ctx context.Context, guildId string, shouldLog bool) error {
	settings, err := m.repo.GetLedgerSettings(ctx, guildId)
	if err != nil {
		return err
	}

	settings.Enabled = shouldLog
	return m.repo.UpdateLedgerSettings(ctx, settings)
}

func (m *RepoLedgerManager) GetLogChannel(ctx context.Context, guildId string) (string, error) {
	settings, err := m.repo.GetLedgerSettings(ctx, guildId)
	if err != nil {
		return "", err
	}

	return settings.LogChannelId, nil
}

func (m *RepoLedgerManager) SetLogChannel(ctx context.Context, guildId string, channelId string) error {
	settings, err := m.repo.GetLedgerSettings(ctx, guildId)
	if err != nil {
		return err
	}

	settings.LogChannelId = channelId
	return m.repo.UpdateLedgerSettings(ctx, settings)
}

func (m *RepoLedgerManager) LogCustomEvent(ctx context.Context, guildId string, data *discordgo.MessageSend) error {
	// Get log channel
	channelId, err := m.GetLogChannel(ctx, guildId)
	if err != nil {
		return err
	}

	_, err = m.session.ChannelMessageSendComplex(channelId, data)
	return err
}

func (m *RepoLedgerManager) LogMessageCreate(ctx context.Context, message *discordgo.Message) error {
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

	// Save in database
	if err = m.repo.CreateMessage(ctx, &LedgerMessage{
		MessageId: message.ID,
		GuildId:   message.GuildID,
		ChannelId: message.ChannelID,
		UserId:    message.Author.ID,
	}); err != nil {
		return err
	}

	// Log contents
	if err = m.repo.CreateMessageContent(ctx, &LedgerContent{
		MessageId: message.ID,
		Content:   message.Content,
	}); err != nil {
		return err
	}

	// Logging creates isn't necessary
	return nil
}

func (m *RepoLedgerManager) LogMessageDelete(ctx context.Context, message *discordgo.Message) error {
	if message.GuildID == "" {
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

	// Save in database
	if err = m.repo.UpdateMessage(ctx, &LedgerMessage{
		MessageId: message.ID,
		IsDeleted: true,
	}); err != nil {
		return err
	}

	messageData, err := m.repo.GetMessage(ctx, message.ID)
	if err != nil {
		return err
	}

	// Check if the message has a user, channel, and guild
	if messageData.UserId == "" || messageData.ChannelId == "" || messageData.GuildId == "" {
		return nil
	}

	// Get author
	author, err := m.session.User(messageData.UserId)
	if err != nil {
		return err
	}

	// Log contents
	content, err := m.repo.GetMessageContents(ctx, message.ID)
	if err != nil {
		return err
	}

	messageContent := "`no content saved in database`"
	if len(content) > 0 {
		messageContent = content[len(content)-1].Content
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
		Color: core.ColorError,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Content",
				Value: messageContent,
			},
			{
				Name:  "URL",
				Value: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", messageData.GuildId, messageData.ChannelId, messageData.UserId),
			},
			{
				Name:   "Author Mention",
				Value:  author.Mention(),
				Inline: true,
			},
			{
				Name:   "Author ID/Tag",
				Value:  fmt.Sprintf("`%s` /\n `%s`", author.ID, author.String()),
				Inline: true,
			},
			{
				Name:   "Channel",
				Value:  "<#" + messageData.ChannelId + ">",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s • %s", message.ID, messageData.CreatedAt.Format("01/02/2006, 3:04:05 PM")),
		},
	}

	_, err = m.session.ChannelMessageSendEmbed(channelId, embed)
	return err
}

func (m *RepoLedgerManager) LogMessageEdit(ctx context.Context, message *discordgo.MessageUpdate) error {
	if message.Author.Bot || message.GuildID == "" {
		return nil
	}

	if message.Message == nil {
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

	// get old message content
	contents, err := m.repo.GetMessageContents(ctx, message.Message.ID)
	if err != nil {
		return err
	}

	previous := "`no content saved in database`"
	if len(contents) > 0 {
		previous = contents[len(contents)-1].Content
	}

	// Save in database
	if err = m.repo.UpdateMessage(ctx, &LedgerMessage{
		MessageId: message.Message.ID,
		GuildId:   message.GuildID,
		ChannelId: message.Message.ChannelID,
		UserId:    message.Message.Author.ID,
		IsEdited:  true,
	}); err != nil {
		return err
	}

	// Log contents
	if err = m.repo.CreateMessageContent(ctx, &LedgerContent{
		MessageId: message.Message.ID,
		Content:   message.Message.Content,
	}); err != nil {
		return err
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
		Title: "Message Edited",
		Color: core.ColorInfo,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Content",
				Value: message.Content,
			},
			{
				Name:  "Previous Content",
				Value: previous,
			},
			{
				Name:  "URL",
				Value: fmt.Sprintf("https://discord.com/channels/%s/%s/%s", message.GuildID, message.ChannelID, message.ID),
			},
			{
				Name:   "Author Mention",
				Value:  message.Author.Mention(),
				Inline: true,
			},
			{
				Name:   "Author ID/Tag",
				Value:  fmt.Sprintf("`%s` /\n `%s`", message.Author.ID, message.Author.String()),
				Inline: true,
			},
			{
				Name:   "Channel",
				Value:  "<#" + message.ChannelID + ">",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%s • %s", message.ID, message.Timestamp.Format("01/02/2006, 3:04:05 PM")),
		},
	}

	_, err = m.session.ChannelMessageSendEmbed(channelId, embed)
	return err
}
