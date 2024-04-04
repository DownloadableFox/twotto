package debug

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

func RegisterModule(client *discordgo.Session) {
	// Add handlers
	onReadyIdent := core.NewIdentifier("debug", "events/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	pingCommandIdent := core.NewIdentifier("debug", "commands/ping")
	pingCommand := core.ApplyMiddlewares(
		HandlePingCommand,
		MidwareForCommand(PingCommand),
		MidwareErrorWrap(pingCommandIdent),
	)
	client.AddHandler(core.HandleEvent(pingCommand))

	errorTestCommandIdent := core.NewIdentifier("debug", "commands/error-test")
	errorTestCommand := core.ApplyMiddlewares(
		HandleErrorTestCommand,
		MidwareForCommand(ErrorTestCommand),
		MidwareErrorWrap(errorTestCommandIdent),
	)
	client.AddHandler(core.HandleEvent(errorTestCommand))
}
