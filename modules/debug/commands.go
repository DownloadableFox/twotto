package debug

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

var ErrorTestCommandPermissions int64 = discordgo.PermissionAdministrator

var ErrorTestCommand = &discordgo.ApplicationCommand{
	Name:                     "error-test",
	Description:              "Development command for testing error handling",
	DefaultMemberPermissions: &ErrorTestCommandPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "no-reply",
			Description: "Throws error before sending a reply.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "reply",
			Description: "Throws error after replying to interaction.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ephemeral",
					Description: "Whether or not the reply should be ephemeral.",
					Type:        discordgo.ApplicationCommandOptionBoolean,
				},
			},
		},
		{
			Name:        "defered",
			Description: "Defer a response before throwing an error.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ephemeral",
					Description: "Whether or not the defer should be ephemeral.",
					Type:        discordgo.ApplicationCommandOptionBoolean,
				},
			},
		},
		{
			Name:        "panic",
			Description: "This will generate a panic in the bot, this option will not reply an error.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
	},
}

var (
	_ core.EventFunc[discordgo.InteractionCreate] = HandleErrorTestCommand
)

func HandleErrorTestCommand(_ context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	data := e.ApplicationCommandData()

	options := data.Options
	switch options[0].Name {
	case "reply":
		var flags discordgo.MessageFlags
		if len(options[0].Options) == 0 || options[0].Options[0].BoolValue() {
			flags |= discordgo.MessageFlagsEphemeral
		}

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: flags,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Meow! :3",
						Color:       core.ColorResult,
						Description: "This is a funny & quirky response! Totally not going to die in the next 2 nanoseconds. An error is about to occur after this, depending on the handling something might or not happen.",
					},
				},
			},
		}

		if err := s.InteractionRespond(e.Interaction, response); err != nil {
			return err
		}
	case "defered":
		var flags discordgo.MessageFlags
		if len(options[0].Options) == 0 || options[0].Options[0].BoolValue() {
			flags |= discordgo.MessageFlagsEphemeral
		}

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{},
		}

		if err := s.InteractionRespond(e.Interaction, response); err != nil {
			return err
		}
	case "panic":
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Welp this hurts!",
						Color:       core.ColorResult,
						Description: "A panic is going to happen in my runtime in the next instants. Please beware that if unhandled correctly this might make me despawn (exit on failure) which wouldn't be optimal.",
					},
				},
			},
		}

		if err := s.InteractionRespond(e.Interaction, response); err != nil {
			return err
		}

		panic("This is a fake panic! Comming from error test command.")
	}

	return errors.New("this is a made up error")
}

var PingCommand = &discordgo.ApplicationCommand{
	Name:        "ping",
	Description: "Ping the bot to see if it's alive!",
}

var (
	_ core.EventFunc[discordgo.InteractionCreate] = HandlePingCommand
)

func HandlePingCommand(_ context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Pong! :3",
					Color:       core.ColorInfo,
					Description: fmt.Sprintf("I am alive and well! Server time is <t:%d:f>.", time.Now().Unix()),
				},
			},
		},
	}

	if err := s.InteractionRespond(e.Interaction, response); err != nil {
		return err
	}

	return nil
}
