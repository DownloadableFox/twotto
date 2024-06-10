package extra

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
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
		Color:       core.ColorSuccess,
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}

var CreateForumCommandDMPermission = false
var CreateForumCommandPermissions int64 = discordgo.PermissionAdministrator

var CreateForumCommand = &discordgo.ApplicationCommand{
	Name:                     "create-forum",
	Description:              "Creates a forum channel.",
	DMPermission:             &CreateForumCommandDMPermission,
	DefaultMemberPermissions: &CreateForumCommandPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "name",
			Description: "The name of the forum.",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func HandleCreateForumCommand(_ context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
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

	// Creates the channel
	channel, err := s.GuildChannelCreate(e.GuildID, data.Options[0].StringValue(), discordgo.ChannelTypeGuildForum)
	if err != nil {
		return err
	}

	// Responds to the interaction
	embed := &discordgo.MessageEmbed{
		Title:       "Forum Created",
		Description: fmt.Sprintf("Forum created at <#%s>.", channel.ID),
		Color:       core.ColorSuccess,
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	}); err != nil {
		return err
	}

	return nil
}
