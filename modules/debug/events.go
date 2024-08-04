package debug

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
)

var (
	_ core.EventFunc[discordgo.MessageCreate] = HandlePongEvent
)

func HandlePongEvent(_ context.Context, s *discordgo.Session, e *discordgo.MessageCreate) error {
	if e.Author.ID == s.State.User.ID {
		return nil
	}

	if strings.Contains(e.Content, "!ping") {
		if _, err := s.ChannelMessageSendReply(e.ChannelID, "pong", e.Reference()); err != nil {
			return err
		}
	}

	return nil
}

func HandleOnReadyEvent(_ context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	log.Info().Msgf("[DebugModule] Logged in as %s#%s", e.User.Username, e.User.Discriminator)

	// Liste guilds
	guilds, err := s.UserGuilds(0, "", "", false)
	if err != nil {
		log.Warn().Err(err).Msg("[DebugModule] Failed to list guilds")
	} else {
		if len(guilds) == 0 {
			log.Warn().Msg("[DebugModule] Not connected to any guilds")
		}

		guildsGreeting := ""
		for i, guild := range guilds {
			if i > 0 {
				guildsGreeting += ", "
			}

			guildsGreeting += guild.Name
		}

		log.Info().Msgf("[DebugModule] Connected to guilds: %s", guildsGreeting)
	}

	// Register commands
	err = core.ApplyCommands(
		PingCommand,
		ErrorTestCommand,
	).For(s, "")
	if err != nil {
		log.Warn().Err(err).Msg("[DebugModule] Failed to register commands")
		return err
	}

	// Update status to do not disturb
	if err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "dnd",
		Activities: []*discordgo.Activity{{
			Name: "furry femboys",
			Type: discordgo.ActivityTypeWatching,
			URL:  "https://www.youtube.com/watch?v=lmSgyD5Jb_w",
		}},
	}); err != nil {
		return err
	}

	return nil
}

func HandleFeatureSetupEvent(c context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	log.Info().Msg("[FeatureServiceSetup] Registering features...")

	fs, ok := c.Value(FeatureServiceKey).(FeatureService)
	if !ok {
		return ErrFeatureServiceNotFound
	}

	features, err := fs.ListFeatures()
	if err != nil {
		log.Warn().Err(err).Msg("[FeatureServiceSetup] Failed to list features!")
		return err
	}

	guilds, err := s.UserGuilds(0, "", "", false)
	if err != nil {
		log.Warn().Err(err).Msg("[FeatureServiceSetup] Failed to list guilds!")
		return err
	}

	for _, guild := range guilds {
		for _, feature := range features {
			_, err := fs.GetFeature(context.Background(), feature.Identifier, guild.ID)
			if err != nil {
				if err == ErrFeatureNotRegistered {
					if err := fs.SetFeature(context.Background(), feature.Identifier, guild.ID, feature.DefaultState); err != nil {
						log.Warn().Err(err).Msgf("[FeatureServiceSetup] Failed to set default feature state for guild \"%s\" (%s)!", guild.Name, guild.ID)
					}

					log.Info().Msgf("[FeatureServiceSetup] Set default feature state for guild \"%s\" (%s)!", guild.Name, guild.ID)
				} else {
					log.Warn().Err(err).Msgf("[FeatureServiceSetup] Failed to get feature state for guild \"%s\" (%s)!", guild.Name, guild.ID)
				}
			}
		}
	}

	log.Info().Msg("[FeatureServiceSetup] Features registered!")

	return nil
}
