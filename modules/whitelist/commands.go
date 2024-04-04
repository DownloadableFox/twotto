package whitelist

import (
	"context"
	"errors"

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
					Name:        "user",
					Description: "The user to add to the whitelist.",
					Type:        discordgo.ApplicationCommandOptionUser,
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
					Name:        "user",
					Description: "The user to remove from the whitelist.",
					Type:        discordgo.ApplicationCommandOptionUser,
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
			Type:        discordgo.ApplicationCommandOptionSubCommand,
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
	return errors.New("not implemented")
}
