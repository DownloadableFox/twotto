package ledger

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var LedgerCommandPermission int64 = discordgo.PermissionAdministrator

var LedgerCommand = &discordgo.ApplicationCommand{
	Name:        "ledger",
	Description: "Manage the ledger module",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "enable",
			Description: "Enable the ledger module",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to log messages to",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "disable",
			Description: "Disable the ledger module",
		},
	},
}

func HandleLedgerCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Defer interaction response
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		return err
	}

	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok || lm == nil {
		return ErrLedgerManagerNotFound
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "enable":
		return HandleEnableLedgerCommand(ctx, s, i)
	case "disable":
		return HandleDisableLedgerCommand(ctx, s, i)
	}
	return nil
}

func HandleEnableLedgerCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	channel := i.ApplicationCommandData().Options[0].Options[0].ChannelValue(s)

	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok || lm == nil {
		return ErrLedgerManagerNotFound
	}

	err := lm.SetLogChannel(ctx, i.GuildID, channel.ID)
	if err != nil {
		return err
	}

	err = lm.SetShouldLog(ctx, i.GuildID, true)
	if err != nil {
		return err
	}

	log.Debug().Msgf("[LedgerModule] Enabled ledger for guild %s, logging to channel %s", i.GuildID, channel.ID)

	// Respond to interaction
	embed := &discordgo.MessageEmbed{
		Title:       "Ledger enabled",
		Color:       0x003DFF,
		Description: "The ledger module has been enabled and messages will now be logged to <#" + channel.ID + ">",
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})

	return err
}

func HandleDisableLedgerCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok || lm == nil {
		return ErrLedgerManagerNotFound
	}

	err := lm.SetShouldLog(ctx, i.GuildID, false)
	if err != nil {
		return err
	}

	log.Debug().Msgf("[LedgerModule] Disabled ledger for guild %s", i.GuildID)

	// Respond to interaction
	embed := &discordgo.MessageEmbed{
		Title:       "Ledger disabled",
		Color:       0x003DFF,
		Description: "The ledger module has been disabled and messages will no longer be logged",
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})

	return err
}
