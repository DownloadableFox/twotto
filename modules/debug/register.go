package debug

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

func RegisterModule(client *discordgo.Session, featureService FeatureService) {
	// Add handlers
	onReadyIdent := core.NewIdentifier("debug", "events/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	featureSetupEventIdent := core.NewIdentifier("debug", "events/feature-setup")
	featureSetupEvent := core.ApplyMiddlewares(
		HandleFeatureSetupEvent,
		MidwareContextInject[discordgo.Ready](FeatureServiceKey, featureService),
		MidwarePerformance[discordgo.Ready](featureSetupEventIdent),
	)
	client.AddHandler(core.HandleEvent(featureSetupEvent))

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
