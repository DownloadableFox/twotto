package whitelist

import (
	"context"
	"errors"

	"github.com/downloadablefox/twotto/core"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	WhitelistManagerKey         = core.NewIdentifier("whitelist", "service/manager")
	ErrWhitelisted              = errors.New("user already whitelisted")
	ErrNotWhitelisted           = errors.New("user not whitelisted")
	ErrNotInGuild               = errors.New("user not in guild")
	ErrDefaultRoleNotFound      = errors.New("default role not set")
	ErrWhitelistManagerNotFound = errors.New("whitelist manager not found in context (missing injection)")
)

type WhitelistManager interface {
	Whitelist(ctx context.Context, guildId string, userId string) error
	Unwhitelist(ctx context.Context, guildId string, userId string) error
	IsWhitelisted(ctx context.Context, guildId string, userId string) bool
	GetWhitelist(ctx context.Context, guildId string) ([]string, error)
	ClearWhitelist(ctx context.Context, guildId string) error
	GetDefaultRole(ctx context.Context, guildId string) string
	SetDefaultRole(ctx context.Context, guildId string, roleId string) error
	GetEnabled(ctx context.Context, guildId string) bool
	SetEnabled(ctx context.Context, guildId string, enabled bool) error
	GetRemoveOnBan(ctx context.Context, guildId string) bool
	SetRemoveOnBan(ctx context.Context, guildId string, removeOnBan bool) error
}

type PostgresWhitelistManager struct {
	pool *pgxpool.Pool
}

func NewPostgresWhitelistManager(pool *pgxpool.Pool) WhitelistManager {
	return &PostgresWhitelistManager{pool: pool}
}

func (m *PostgresWhitelistManager) IsWhitelisted(ctx context.Context, guildId string, userId string) bool {
	var whitelisted bool
	err := m.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM whitelist_users
			WHERE guild_id = $1 AND user_id = $2
		)
	`, guildId, userId).Scan(&whitelisted)
	if err != nil {
		return false
	}

	return whitelisted
}

func (m *PostgresWhitelistManager) Whitelist(ctx context.Context, guildId string, userId string) error {
	if m.IsWhitelisted(ctx, guildId, userId) {
		return ErrWhitelisted
	}

	_, err := m.pool.Exec(ctx, `
		INSERT INTO whitelist_users (guild_id, user_id)
		VALUES ($1, $2)
	`, guildId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresWhitelistManager) Unwhitelist(ctx context.Context, guildId string, userId string) error {
	if !m.IsWhitelisted(ctx, guildId, userId) {
		return ErrNotWhitelisted
	}

	_, err := m.pool.Exec(ctx, `
		DELETE FROM whitelist_users
		WHERE guild_id = $1 AND user_id = $2
	`, guildId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresWhitelistManager) GetWhitelist(ctx context.Context, guildId string) ([]string, error) {
	rows, err := m.pool.Query(ctx, `
		SELECT user_id
		FROM whitelist_users
		WHERE guild_id = $1
	`, guildId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var whitelist []string
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}

		whitelist = append(whitelist, userId)
	}

	return whitelist, nil
}

func (m *PostgresWhitelistManager) ClearWhitelist(ctx context.Context, guildId string) error {
	_, err := m.pool.Exec(ctx, `
		DELETE FROM whitelist_users
		WHERE guild_id = $1
	`, guildId)

	return err
}

func (m *PostgresWhitelistManager) GetDefaultRole(ctx context.Context, guildId string) string {
	var roleId string
	err := m.pool.QueryRow(ctx, `
		SELECT default_role_id
		FROM whitelist_settings
		WHERE guild_id = $1
	`, guildId).Scan(&roleId)
	if err != nil {
		return ""
	}

	return roleId
}

func (m *PostgresWhitelistManager) SetDefaultRole(ctx context.Context, guildId string, roleId string) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO whitelist_settings (guild_id, default_role_id)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE SET default_role_id = $2
	`, guildId, roleId)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresWhitelistManager) GetEnabled(ctx context.Context, guildId string) bool {
	var enabled bool
	err := m.pool.QueryRow(ctx, `
		SELECT enabled
		FROM whitelist_settings
		WHERE guild_id = $1
	`, guildId).Scan(&enabled)
	if err != nil {
		return false
	}

	return enabled
}

func (m *PostgresWhitelistManager) SetEnabled(ctx context.Context, guildId string, enabled bool) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO whitelist_settings (guild_id, enabled)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE SET enabled = $2
	`, guildId, enabled)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresWhitelistManager) GetRemoveOnBan(ctx context.Context, guildId string) bool {
	var removeOnBan bool
	err := m.pool.QueryRow(ctx, `
		SELECT remove_on_ban
		FROM whitelist_settings
		WHERE guild_id = $1
	`, guildId).Scan(&removeOnBan)
	if err != nil {
		return false
	}

	return removeOnBan
}

func (m *PostgresWhitelistManager) SetRemoveOnBan(ctx context.Context, guildId string, removeOnBan bool) error {
	_, err := m.pool.Exec(ctx, `
		INSERT INTO whitelist_settings (guild_id, remove_on_ban)
		VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE SET remove_on_ban = $2
	`, guildId, removeOnBan)
	if err != nil {
		return err
	}

	return nil
}
