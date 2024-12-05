package debug

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
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

var FeatureCommandPermissions int64 = discordgo.PermissionAdministrator

var FeatureCommand = &discordgo.ApplicationCommand{
	Name:                     "feature",
	Description:              "Manage features for the bot.",
	DefaultMemberPermissions: &FeatureCommandPermissions,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "get",
			Description: "Get the state of a feature.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "feature",
					Description:  "The feature to get the state of.",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "set",
			Description: "Set the state of a feature.",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "feature",
					Description:  "The feature to set the state of.",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
				{
					Name:        "state",
					Description: "The state to set the feature to.",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Required:    true,
				},
			},
		},
	},
}

var (
	_ core.EventFunc[discordgo.InteractionCreate] = HandleFeatureCommand
	_ core.EventFunc[discordgo.InteractionCreate] = HandleFeatureAutocomplete
)

func HandleFeatureCommand(c context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	data := e.ApplicationCommandData()

	fs, ok := c.Value(FeatureServiceKey).(FeatureService)
	if !ok {
		return ErrFeatureServiceNotFound
	}

	options := data.Options
	switch options[0].Name {
	case "get":
		featureName := options[0].Options[0].StringValue()
		identifier, err := core.ParseIdentifier(featureName)
		if err != nil {
			return err
		}

		enabled, err := fs.GetFeature(context.Background(), identifier, e.GuildID)
		if err != nil {
			if errors.Is(err, ErrFeatureNotRegistered) {
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags: discordgo.MessageFlagsEphemeral,
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "Feature not registered!",
								Color:       core.ColorError,
								Description: fmt.Sprintf("The feature `%s` is not registered for this guild.", featureName),
							},
						},
					},
				}

				if err := s.InteractionRespond(e.Interaction, response); err != nil {
					return err
				}

				return nil
			}

			return err
		}

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Feature state",
						Color:       core.ColorInfo,
						Description: fmt.Sprintf("The feature `%s` is currently %s.", featureName, map[bool]string{true: "enabled", false: "disabled"}[enabled]),
					},
				},
			},
		}

		if err := s.InteractionRespond(e.Interaction, response); err != nil {
			return err
		}
	case "set":
		featureName := options[0].Options[0].StringValue()
		identifier, err := core.ParseIdentifier(featureName)
		if err != nil {
			return err
		}

		state := options[0].Options[1].BoolValue()

		if err := fs.SetFeature(context.Background(), identifier, e.GuildID, state); err != nil {
			return err
		}

		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Feature state updated!",
						Color:       core.ColorSuccess,
						Description: fmt.Sprintf("The feature `%s` is now %s.", featureName, map[bool]string{true: "enabled", false: "disabled"}[state]),
					},
				},
			},
		}

		if err := s.InteractionRespond(e.Interaction, response); err != nil {
			return err
		}
	}

	return nil
}

func HandleFeatureAutocomplete(c context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	fs, ok := c.Value(FeatureServiceKey).(FeatureService)
	if !ok {
		return ErrFeatureServiceNotFound
	}

	features, err := fs.ListFeatures()
	if err != nil {
		return err
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(features))
	for _, feature := range features {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  feature.Identifier.String(),
			Value: feature.Identifier.String(),
		})
	}

	if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}); err != nil {
		return err
	}

	return nil
}

var RestartCommand = &discordgo.ApplicationCommand{
	Name:        "restart",
	Description: "Restarts the bot.",
	Version:     "1.0.0",
}

func HandleRestartCommand(c context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	var owners = []string{"556132236697665547", "836684190987583576", "610825796285890581"}

	var userId string
	if e.Member != nil {
		userId = e.Member.User.ID
	} else {
		userId = e.User.ID
	}

	if !slices.Contains(owners, userId) {
		return errors.New("you are not authorized to use this command")
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Restarting...",
					Color:       core.ColorInfo,
					Description: "The bot is now restarting. Please wait a moment.",
				},
			},
		},
	}

	if err := s.InteractionRespond(e.Interaction, response); err != nil {
		return err
	}

	// Log the restart
	log.Info().Msg("[Debug] Restart command received- Restarting bot...")

	// Remove all commands
	core.UnregisterAllCommands(s)

	// Close the session
	os.Exit(1)

	return nil
}
