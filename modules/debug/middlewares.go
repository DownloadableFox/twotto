package debug

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

func MidwarePerformance[T any](tag *core.Identifier) core.MiddlewareFunc[T] {
	return func(next core.EventFunc[T]) core.EventFunc[T] {
		return func(c context.Context, s *discordgo.Session, e *T) error {
			start := time.Now()
			err := next(c, s, e)
			log.Info().Msgf("[PerformanceMidware] Finished event execution for \"%s\", took: %s", tag, time.Since(start))

			return err
		}
	}
}

func CreateErrorEmbed(err error, id xid.ID) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xFF003D,
		Title:       "Oh no! :(",
		Description: "Sorry! An unexpected error occurred while executing this event.\nIf this keeps happening contact <@556132236697665547>.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Error Message",
				Value: err.Error(),
			},
			{
				Name:   "Error ID",
				Value:  fmt.Sprintf("`%s`", id),
				Inline: true,
			},
			{
				Name:   "Server Time",
				Value:  fmt.Sprintf("<t:%d:f>", id.Time().Unix()),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "The following error was reported!",
		},
	}
}

func CreateFatalErrorEmbed(id xid.ID) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       0xFF003D,
		Title:       "Fatal! -w-",
		Description: "You have encountered a fatal error! This should never happen.\nIf this keeps happening contact <@556132236697665547>.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Error ID",
				Value:  fmt.Sprintf("`%s`", id),
				Inline: true,
			},
			{
				Name:   "Server Time",
				Value:  fmt.Sprintf("<t:%d:f>", id.Time().Unix()),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "The following error was reported!",
		},
	}
}

func panicWrap(s *discordgo.Session, e *discordgo.InteractionCreate) {
	rec := recover()
	if rec == nil {
		return
	}

	log.Error().Any("panic", rec).Msg("[ErrorWrapMidware] Recovered from panic!")

	// Generate ID
	id := xid.New()

	// Generate embed
	errorEmbed := CreateFatalErrorEmbed(id)

	// Attempt to reply
	if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{errorEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}); err == nil { // If reply succesful finish.
		return
	}

	// Check for reply
	res, err := s.InteractionResponse(e.Interaction)
	if err != nil {
		return
	}

	if res.Flags&discordgo.MessageFlagsLoading > 0 {
		_, err = s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{errorEmbed},
		})
		if err != nil {
			return
		}
	}

	// Get stacktrace
	stacktrace := make([]byte, 4096)
	count := runtime.Stack(stacktrace, false)

	// Create reader
	reader := bytes.NewReader(stacktrace[:count])

	s.FollowupMessageCreate(e.Interaction, false, &discordgo.WebhookParams{
		Flags: res.Flags & discordgo.MessageFlagsEphemeral,
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("st-%s.txt", id),
				ContentType: "text/plain",
				Reader:      reader,
			},
		},
	})

}

func MidwareErrorWrap(tag *core.Identifier) core.MiddlewareFunc[discordgo.InteractionCreate] {
	return func(next core.EventFunc[discordgo.InteractionCreate]) core.EventFunc[discordgo.InteractionCreate] {
		return func(c context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
			defer panicWrap(s, e)

			if err := next(c, s, e); err != nil {
				// Generate ID
				id := xid.New()

				// Generate embed
				errorEmbed := CreateErrorEmbed(err, id)

				log.Warn().Err(err).Msgf("[ErrorWrapMidware] Caught an error while executing interaction \"%s\"!", tag)

				// Attempt to reply
				if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{errorEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				}); err == nil { // If reply succesful finish.
					return nil
				}

				// Check for reply
				res, err := s.InteractionResponse(e.Interaction)
				if err != nil {
					// Failed to respond
					return fmt.Errorf("failed to respond to interaction \"%s\" with error: \"%s\"", tag, err.Error())
				}

				// Edit reply where possible
				if res.Flags&discordgo.MessageFlagsLoading > 0 {
					_, err = s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
						Embeds: &[]*discordgo.MessageEmbed{errorEmbed},
					})
					return err
				}
			}

			return nil
		}
	}
}

func MidwareForCommand(command *discordgo.ApplicationCommand) core.MiddlewareFunc[discordgo.InteractionCreate] {
	// Register command
	return func(next core.EventFunc[discordgo.InteractionCreate]) core.EventFunc[discordgo.InteractionCreate] {
		return func(c context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
			data := e.ApplicationCommandData()

			// Ignore if not the correct command
			if data.Name != command.Name {
				return nil
			}

			return next(c, s, e)
		}
	}
}

func MidwareContextInject[T interface{}](key *core.Identifier, value any) core.MiddlewareFunc[T] {
	return func(next core.EventFunc[T]) core.EventFunc[T] {
		return func(c context.Context, s *discordgo.Session, e *T) error {
			c = context.WithValue(c, key, value)
			return next(c, s, e)
		}
	}
}

func GetGuildFromEvent(event interface{}) string {
	switch e := event.(type) {
	case *discordgo.InteractionCreate:
		return e.GuildID
	case *discordgo.MessageCreate:
		return e.GuildID
	default:
		return ""
	}
}

func MidwareModuleEnabled[T interface{}](module string, provider ModuleProvider) core.MiddlewareFunc[T] {
	return func(next core.EventFunc[T]) core.EventFunc[T] {
		return func(c context.Context, s *discordgo.Session, e *T) error {
			guildId := GetGuildFromEvent(e)
			if guildId == "" {
				return fmt.Errorf("failed to determine guild id for event %T", e)
			}
			if !provider.IsModuleEnabled(module, guildId) {
				return nil
			}

			return next(c, s, e)
		}
	}
}

func MidwareFeatureEnabled[T interface{}](identifier *core.Identifier, provider ModuleProvider) core.MiddlewareFunc[T] {
	return func(next core.EventFunc[T]) core.EventFunc[T] {
		return func(c context.Context, s *discordgo.Session, e *T) error {
			guildId := GetGuildFromEvent(e)
			if guildId == "" {
				return fmt.Errorf("failed to determine guild id for event %T", e)
			}
			if !provider.IsFeatureEnabled(identifier, guildId) {
				return nil
			}

			return next(c, s, e)
		}
	}
}
