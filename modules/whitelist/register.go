package whitelist

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
)

func RegisterModule(client *discordgo.Session) {
	onReadyIdent := core.NewIdentifier("whitelist", "events/setup")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwareContextInject[discordgo.Ready](WhitelistManagerKey, nil),
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	onJoin := core.ApplyMiddlewares(
		HandleOnJoinEvent,
		debug.MidwareContextInject[discordgo.GuildMemberAdd](WhitelistManagerKey, nil),
	)
	client.AddHandler(core.HandleEvent(onJoin))

	onBan := core.ApplyMiddlewares(
		HandleOnBanEvent,
		debug.MidwareContextInject[discordgo.GuildBanAdd](WhitelistManagerKey, nil),
	)
	client.AddHandler(core.HandleEvent(onBan))

	whitelistCommandIndent := core.NewIdentifier("whitelist", "commands/whitelist")
	whitelistCommand := core.ApplyMiddlewares(
		HandleWhitelistCommand,
		debug.MidwareContextInject[discordgo.InteractionCreate](WhitelistManagerKey, nil),
		debug.MidwareForCommand(WhitelistCommand),
		debug.MidwareErrorWrap(whitelistCommandIndent),
	)
	client.AddHandler(core.HandleEvent(whitelistCommand))
}
