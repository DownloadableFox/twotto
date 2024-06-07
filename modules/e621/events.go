package e621

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
)

func HandleOnReadyEvent(_ context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	// Register commands
	err := core.ApplyCommands(
		YiffCommand,
	).For(s, "")
	if err != nil {
		log.Warn().Err(err).Msg("[DebugModule] Failed to register commands")
		return err
	}

	return nil
}
