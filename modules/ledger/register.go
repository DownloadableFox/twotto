package ledger

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/downloadablefox/twotto/modules/debug"
)

func RegisterModule(client *discordgo.Session, ledger LedgerManager) {
	createMessage := core.ApplyMiddlewares(
		HandleOnMessageCreateEvent,
		debug.MidwareContextInject[discordgo.MessageCreate](LedgerManagerKey, ledger),
	)
	client.AddHandler(core.HandleEvent(createMessage))

	editMessage := core.ApplyMiddlewares(
		HandleOnMessageEditEvent,
		debug.MidwareContextInject[discordgo.MessageUpdate](LedgerManagerKey, ledger),
	)
	client.AddHandler(core.HandleEvent(editMessage))

	deleteMessage := core.ApplyMiddlewares(
		HandleOnMessageDeleteEvent,
		debug.MidwareContextInject[discordgo.MessageDelete](LedgerManagerKey, ledger),
	)
	client.AddHandler(core.HandleEvent(deleteMessage))

	onReadyIdent := core.NewIdentifier("ledger", "event/ready")
	onReady := core.ApplyMiddlewares(
		HandleOnReadyEvent,
		debug.MidwareContextInject[discordgo.Ready](LedgerManagerKey, ledger),
		debug.MidwarePerformance[discordgo.Ready](onReadyIdent),
	)
	client.AddHandler(core.HandleEvent(onReady))

	ledgerCommandIdent := core.NewIdentifier("ledger", "commands/ledger")
	ledgerCommand := core.ApplyMiddlewares(
		HandleLedgerCommand,
		debug.MidwareContextInject[discordgo.InteractionCreate](LedgerManagerKey, ledger),
		debug.MidwareForCommand(LedgerCommand),
		debug.MidwareErrorWrap(ledgerCommandIdent),
	)
	client.AddHandler(core.HandleEvent(ledgerCommand))
}
