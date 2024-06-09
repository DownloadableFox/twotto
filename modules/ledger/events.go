package ledger

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

func HandleOnMessageCreateEvent(ctx context.Context, s *discordgo.Session, e *discordgo.MessageCreate) error {
	// Ignore bots
	if e.Author != nil && e.Author.Bot {
		return nil
	}

	// Ignore messages from whitelisted users
	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok {
		return ErrLedgerManagerNotFound
	}

	// Log message
	return lm.LogMessageCreate(ctx, e.Message)
}

func HandleOnMessageEditEvent(ctx context.Context, s *discordgo.Session, e *discordgo.MessageUpdate) error {
	// Ignore bots
	if e.Author != nil && e.Author.Bot {
		return nil
	}

	// Ignore messages from whitelisted users
	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok {
		return ErrLedgerManagerNotFound
	}

	// Log message
	return lm.LogMessageEdit(ctx, e)
}

func HandleOnMessageDeleteEvent(ctx context.Context, s *discordgo.Session, e *discordgo.MessageDelete) error {
	// Ignore bots
	if e.Author != nil && e.Author.Bot {
		return nil
	}

	// Ignore messages from whitelisted users
	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok {
		return ErrLedgerManagerNotFound
	}

	// Log message
	return lm.LogMessageDelete(ctx, e.Message)
}

func HandleOnReadyEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	// Ignore messages from whitelisted users
	lm, ok := ctx.Value(LedgerManagerKey).(LedgerManager)
	if !ok || lm == nil {
		return ErrLedgerManagerNotFound
	}

	// Register slash command
	if err := core.ApplyCommands(LedgerCommand).For(s, ""); err != nil {
		return err
	}

	return nil
}
