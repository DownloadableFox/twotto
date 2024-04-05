package whitelist

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var WhitelistCommandPermissions int64 = discordgo.PermissionAdministrator

var WhitelistCommand = &discordgo.ApplicationCommand{
	Name:                     "whitelist",
	Description:              "Manage the whitelist for the bot.",
	DefaultMemberPermissions: &WhitelistCommandPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "add",
			Description: "Add a user to the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user-id",
					Description: "The id of the user to add to the whitelist.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "remove",
			Description: "Remove a user from the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user-id",
					Description: "The id of the user to remove from the whitelist.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "list",
			Description: "List all users on the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "clear",
			Description: "Clear the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "add-all",
			Description: "Adds all the users in the guild to the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "config",
			Description: "Configure the whitelist.",
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "enable",
					Description: "Enable the whitelist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "disable",
					Description: "Disable the whitelist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "status",
					Description: "Check the status of the whitelist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "set-role",
					Description: "Set the default role for the whitelist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "The role to set as the default role.",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
				{
					Name:        "clear-role",
					Description: "Clear the default role for the whitelist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "set-remove-on-ban",
					Description: "Remove users from the whitelist when they are banned.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "enabled",
							Description: "Whether or not to remove users from the whitelist when they are banned.",
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Required:    true,
						},
					},
				},
			},
		},
	},
}

func HandleWhitelistCommand(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	// Defers the response
	if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		return err
	}

	if ws, ok := ctx.Value(WhitelistManagerKey).(WhitelistManager); !ok || ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Handle the subcommand
	switch e.ApplicationCommandData().Options[0].Name {
	case "add":
		return HandleWhitelistCommandAdd(ctx, s, e)
	case "remove":
		return HandleWhitelistCommandRemove(ctx, s, e)
	case "list":
		return HandleWhitelistCommandList(ctx, s, e)
	case "clear":
		return HandleWhitelistCommandClear(ctx, s, e)
	case "add-all":
		return HandleWhitelistCommandAddAll(ctx, s, e)
	case "config":
		return HandleWhitelistCommandConfig(ctx, s, e)
	default:
		return errors.New("subcommand not yet implemented")
	}
}

func HandleWhitelistCommandAdd(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the user to add
	userId := e.ApplicationCommandData().Options[0].Options[0].StringValue()

	// Add the user to the whitelist
	if err := ws.Whitelist(ctx, e.GuildID, userId); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "User Whitelisted",
		Description: fmt.Sprintf("The user <@%s> has been added to the whitelist.", userId),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandRemove(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the user to remove
	userId := e.ApplicationCommandData().Options[0].Options[0].StringValue()

	// Remove the user from the whitelist
	if err := ws.Unwhitelist(ctx, e.GuildID, userId); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "User Removed from Whitelist",
		Description: fmt.Sprintf("The user <@%s> has been removed from the whitelist.", userId),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandList(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the whitelist
	whitelist, err := ws.GetWhitelist(ctx, e.GuildID)
	if err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Whitelist",
		Description: fmt.Sprintf("There are %d users on the whitelist.", len(whitelist)),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandClear(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Clear the whitelist
	if err := ws.ClearWhitelist(ctx, e.GuildID); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Whitelist Cleared",
		Description: "The whitelist has been cleared.",
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandAddAll(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the members
	members, err := s.GuildMembers(e.GuildID, "", 1000)
	if err != nil {
		return err
	}

	// Whitelist all the members
	for _, member := range members {
		if member.User.Bot {
			continue
		}

		if err := ws.Whitelist(ctx, e.GuildID, member.User.ID); err != nil {
			return err
		}
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "All Users Whitelisted",
		Description: fmt.Sprintf("All %d users in the guild have been added to the whitelist.", len(members)),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfig(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Handle the subcommand
	switch e.ApplicationCommandData().Options[0].Options[0].Name {
	case "enable":
		return HandleWhitelistCommandConfigEnable(ctx, s, e)
	case "disable":
		return HandleWhitelistCommandConfigDisable(ctx, s, e)
	case "status":
		return HandleWhitelistCommandConfigStatus(ctx, s, e)
	case "set-role":
		return HandleWhitelistCommandConfigSetRole(ctx, s, e)
	case "clear-role":
		return HandleWhitelistCommandConfigClearRole(ctx, s, e)
	case "set-remove-on-ban":
		return HandleWhitelistCommandConfigSetRemoveOnBan(ctx, s, e)
	default:
		return errors.New("subcommand not yet implemented")
	}
}

func HandleWhitelistCommandConfigEnable(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Enable the whitelist
	if err := ws.SetEnabled(ctx, e.GuildID, true); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Whitelist Enabled",
		Description: "The whitelist has been enabled.",
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfigDisable(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Disable the whitelist
	if err := ws.SetEnabled(ctx, e.GuildID, false); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Whitelist Disabled",
		Description: "The whitelist has been disabled.",
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfigStatus(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the status
	enabled := ws.GetEnabled(ctx, e.GuildID)

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Whitelist Status",
		Description: fmt.Sprintf("The whitelist is currently %s.", map[bool]string{true: "enabled", false: "disabled"}[enabled]),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfigSetRole(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the role to set
	role := e.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, e.GuildID)

	// Set the default role
	if err := ws.SetDefaultRole(ctx, e.GuildID, role.ID); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Default Role Set",
		Description: fmt.Sprintf("The default role has been set to <@&%s>.", role.ID),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfigClearRole(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Clear the default role
	if err := ws.SetDefaultRole(ctx, e.GuildID, ""); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Default Role Cleared",
		Description: "The default role has been cleared.",
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

func HandleWhitelistCommandConfigSetRemoveOnBan(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	ws := ctx.Value(WhitelistManagerKey).(WhitelistManager)
	if ws == nil {
		return ErrWhitelistManagerNotFound
	}

	// Get the enabled value
	enabled := e.ApplicationCommandData().Options[0].Options[0].Options[0].BoolValue()

	// Set the remove on ban
	if err := ws.SetRemoveOnBan(ctx, e.GuildID, enabled); err != nil {
		return err
	}

	// Respond to the interaction (deferred)
	embed := &discordgo.MessageEmbed{
		Title:       "Remove on Ban Set",
		Description: fmt.Sprintf("Users will now %s be removed from the whitelist when they are banned.", map[bool]string{true: "", false: "not"}[enabled]),
		Color:       0xAE00FF,
	}

	// Edit the response
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}
