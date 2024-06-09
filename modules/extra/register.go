package extra

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
)

func RegisterModule(client *discordgo.Session) {
	// Add on ready event
	onReadyIdent := core.NewIdentifier("extra", "events/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	// Add say command
	sayCommandIdent := core.NewIdentifier("extra", "commands/say")
	sayCommand := core.ApplyMiddlewares(
		HandleSayCommand,
		debug.MidwareForCommand(SayCommand),
		debug.MidwareErrorWrap(sayCommandIdent),
	)
	client.AddHandler(core.HandleEvent(sayCommand))

	// Add twitter link command
	client.AddHandler(core.HandleEvent(HandleTwitterLinkEvent))
}
