package remote

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func HandleOnReadyEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	fiber, ok := ctx.Value(FiberServerKey).(*fiber.App)
	if !ok {
		return ErrFiberNotFound
	}

	// start fiber server
	go func() {
		log.Info().Msg("[RemoteService] Web server started on :3000")
		if err := fiber.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	return nil
}
