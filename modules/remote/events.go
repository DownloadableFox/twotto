package remote

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func HandleOnReadyEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	// Ignore messages from whitelisted users
	web, ok := ctx.Value(WebServerKey).(*fiber.App)
	if !ok || web == nil {
		return ErrWebServerNotFound
	}

	// Start web server
	go func() {
		log.Info().Msg("Starting web server on :3000")

		if err := web.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	return nil
}

func HandleOnShutdownEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Disconnect) error {
	// Ignore messages from whitelisted users
	web, ok := ctx.Value(WebServerKey).(*fiber.App)
	if !ok || web == nil {
		return ErrWebServerNotFound
	}

	// Shutdown web server
	if err := web.Shutdown(); err != nil {
		return err
	}

	return nil
}
