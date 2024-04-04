package whitelist

import (
	"errors"

	"github.com/downloadablefox/twotto/core"
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
	Whitelist(guildId string, userId string) error
	Unwhitelist(guildId string, userId string) error
	IsWhitelisted(guildId string, userId string) bool
	GetDefaultRole(guildId string) string
}
