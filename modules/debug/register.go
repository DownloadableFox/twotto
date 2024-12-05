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

	featureCommandIdent := core.NewIdentifier("debug", "commands/feature")
	featureCommand := core.ApplyMiddlewares(
		HandleFeatureCommand,
		MidwareContextInject[discordgo.InteractionCreate](FeatureServiceKey, featureService),
		MidwareForCommand(FeatureCommand),
		MidwareErrorWrap(featureCommandIdent),
	)
	client.AddHandler(core.HandleEvent(featureCommand))

	featureCommandAutoComplete := core.ApplyMiddlewares(
		HandleFeatureAutocomplete,
		MidwareContextInject[discordgo.InteractionCreate](FeatureServiceKey, featureService),
		MidwareForAutocomplete(FeatureCommand),
	)
	client.AddHandler(core.HandleEvent(featureCommandAutoComplete))

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

	restartCommandIdent := core.NewIdentifier("debug", "commands/restart")
	restartCommand := core.ApplyMiddlewares(
		HandleRestartCommand,
		MidwareForCommand(RestartCommand),
		MidwareErrorWrap(restartCommandIdent),
	)

	client.AddHandler(core.HandleEvent(restartCommand))
}
