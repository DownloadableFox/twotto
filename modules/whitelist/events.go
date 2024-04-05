package whitelist

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
)

func CreateKickInfoEmbed(session *discordgo.Session, userId string, guildId string) (*discordgo.MessageEmbed, error) {
	user, err := session.User(userId)
	if err != nil {
		return nil, errors.New("failed to create kick info embed, user not found")
	}

	guild, err := session.Guild(guildId)
	if err != nil {
		return nil, errors.New("failed to create kick info embed, guild not found")
	}

	return &discordgo.MessageEmbed{
		Title:       "Sorry! :(",
		Description: fmt.Sprintf("You have been kicked from the server `%s` because you are not whitelisted.\nIf this is an error please contact <@556132236697665547> and give her your user ID.", guild.Name),
		Color:       0xff003d,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "User ID",
				Value: fmt.Sprintf("Your user ID is: `%s`", user.ID),
			},
		},
	}, nil
}

func HandleOnJoinEvent(ctx context.Context, s *discordgo.Session, e *discordgo.GuildMemberAdd) error {
	wm, ok := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if !ok {
		return ErrWhitelistManagerNotFound
	}

	// Ignore bots
	if e.User.Bot {
		return nil
	}

	if !wm.IsWhitelisted(ctx, e.GuildID, e.User.ID) {
		log.Warn().Msgf("[WhitelistModule] User %s (%s) joined guild %s but is not whitelisted! Kicking...", e.User, e.User.ID, e.GuildID)

		// Attempt to DM the user
		if dm, err := s.UserChannelCreate(e.User.ID); err == nil {
			if embed, err := CreateKickInfoEmbed(s, e.User.ID, e.GuildID); err != nil {
				log.Warn().Err(err).Msgf("[WhitelistModule] Failed to create kick info embed for user %s (%s)", e.User, e.User.ID)
			} else {
				if _, err := s.ChannelMessageSendEmbed(dm.ID, embed); err != nil {
					log.Warn().Err(err).Msgf("[WhitelistModule] Failed to send kick info to user %s (%s)", e.User, e.User.ID)
				}
			}
		}

		if err := s.GuildMemberDeleteWithReason(e.GuildID, e.User.ID, "Not whitelisted"); err != nil {
			return err
		}

		return nil
	}

	if role := wm.GetDefaultRole(ctx, e.GuildID); role != "" {
		if err := s.GuildMemberRoleAdd(e.GuildID, e.User.ID, role); err != nil {
			return err
		}
	}

	return nil
}

func HandleOnBanEvent(ctx context.Context, s *discordgo.Session, e *discordgo.GuildBanAdd) error {
	wm, ok := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if !ok {
		return ErrWhitelistManagerNotFound
	}

	// Ignore bots
	if e.User.Bot {
		return nil
	}

	if wm.IsWhitelisted(ctx, e.GuildID, e.User.ID) {
		log.Warn().Msgf("[WhitelistModule] User %s (%s) was banned from guild %s, removing user from the whitelist...", e.User, e.User.ID, e.GuildID)

		if err := wm.Unwhitelist(ctx, e.GuildID, e.User.ID); err != nil {
			return err
		}
	}

	return nil
}

func HandleOnReadyEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	_, ok := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if !ok {
		log.Warn().Msg("[WhitelistModule] Whitelist manager not found in context (missing injection), failed to register commands")
		return ErrWhitelistManagerNotFound
	}

	// Register
	if err := core.ApplyCommands(
		WhitelistCommand,
	).For(s, ""); err != nil {
		log.Warn().Msg("[WhitelistModule] Failed to register commands")
		return err
	}

	return nil
}
