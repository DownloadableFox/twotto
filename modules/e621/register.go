package e621

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
)

func RegisterModule(client *discordgo.Session, e621Service IE621Service) {
	// Add handlers
	onReadyIdent := core.NewIdentifier("e621", "events/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	yiffCommandIdent := core.NewIdentifier("e621", "commands/yiff")
	pingCommand := core.ApplyMiddlewares(
		HandleYiffCommand,
		debug.MidwareContextInject[discordgo.InteractionCreate](E621ServiceKey, e621Service),
		debug.MidwareForCommand(YiffCommand),
		debug.MidwareErrorWrap(yiffCommandIdent),
	)
	client.AddHandler(core.HandleEvent(pingCommand))
}
