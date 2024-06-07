package extra

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

var (
	SayCommandDMPermission       = false
	SayCommandPermissions  int64 = discordgo.PermissionAdministrator
)

var SayCommand = &discordgo.ApplicationCommand{
	Name:                     "say",
	Description:              "Say something as the bot!",
	DMPermission:             &SayCommandDMPermission,
	DefaultMemberPermissions: &SayCommandPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "message",
			Description: "The message to say.",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func HandleSayCommand(_ context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	// Defers the response
	if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		return err
	}

	data := e.ApplicationCommandData()

	if _, err := s.ChannelMessageSend(e.ChannelID, data.Options[0].StringValue()); err != nil {
		return err
	}

	// Responds to the interaction
	embed := &discordgo.MessageEmbed{
		Title:       "Message Sent",
		Description: data.Options[0].StringValue(),
		Color:       0x00ff00,
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}
