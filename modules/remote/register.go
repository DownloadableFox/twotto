package remote

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
	"github.com/gofiber/fiber/v2"
)

func RegisterModule(client *discordgo.Session, web *fiber.App) error {
	// Register the remote module
	remote := web.Group("/remote/v1")
	remote.Use(WebMidwareLogger)
	remote.Get("/heartbeat", HandleHeartbeat)
	remote.Get("/activity", HandleGetActivity(client))

	// Register the events
	onReadyIdent := core.NewIdentifier("remote", "event/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwareContextInject[discordgo.Ready](FiberServerKey, web),
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	return nil
}
