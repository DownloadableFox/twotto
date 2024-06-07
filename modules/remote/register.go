package remote

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
	"github.com/gofiber/fiber/v2"
)

var (
	WebServerKey         = core.NewIdentifier("remote", "service/webserver")
	ErrWebServerNotFound = errors.New("web server not found")
)

/*
Disabled for now, may not work correctly. I'd probably change how the system works.
*/
func RegisterModule(client *discordgo.Session, webserver *WebServer) {
	// Add on ping callback
	webserver.AddCallback(func(app *fiber.App) {
		app.Get("/ping", func(c *fiber.Ctx) error {
			return c.SendString("pong")
		})
	})

	onReadyIdent := core.NewIdentifier("remote", "event/setup")

	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwareContextInject[discordgo.Ready](WebServerKey, webserver),
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))
}
