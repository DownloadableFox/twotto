//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/modules/debug"
	"github.com/downloadablefox/twotto/modules/e621"
	"github.com/downloadablefox/twotto/modules/extra"
	"github.com/downloadablefox/twotto/modules/ledger"
	"github.com/downloadablefox/twotto/modules/whitelist"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Token       string `usage:"Discord bot token" required:"true" env:"TOKEN"`
	Debug       bool   `usage:"Enable debug mode" default:"false" env:"DEBUG"`
	DatabaseURL string `usage:"Postgres database URL" required:"true" env:"DATABASE_URL"`
}

func InitializeDatabasePool(config *Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

func InitializeFiberServer() *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var e *fiber.Error
			if errors.As(err, &e) {
				return c.Status(e.Code).JSON(fiber.Map{
					"status":  "error",
					"code":    e.Code,
					"message": e.Message,
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"code":    500,
				"message": err.Error(),
			})
		},
	})
}

func InitializeWhitelistManager(pool *pgxpool.Pool) whitelist.WhitelistManager {
	wire.Build(whitelist.NewPostgresWhitelistManager)
	return nil
}

func InitializeLedgerManager(session *discordgo.Session, pool *pgxpool.Pool) ledger.LedgerManager {
	wire.Build(ledger.NewRepoLedgerManager, ledger.NewLedgerPostgresRepository)
	return nil
}

func InitializeFeatureService(pool *pgxpool.Pool) debug.FeatureService {
	wire.Build(debug.NewPostgresFeatureService)
	return nil
}

func bootstrap(client *discordgo.Session, config *Config) error {
	// Set intents
	client.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentGuildMessageReactions | discordgo.IntentGuildMembers | discordgo.IntentGuildBans
	client.StateEnabled = true
	client.State.TrackMembers = true
	client.State.TrackPresences = true

	pool, err := InitializeDatabasePool(config)
	if err != nil {
		return err
	}

	// Register modules
	featureService := InitializeFeatureService(pool)
	debug.RegisterModule(client, featureService)
	extra.RegisterModule(client, featureService)

	whitelistManager := InitializeWhitelistManager(pool)
	whitelist.RegisterModule(client, whitelistManager)

	ledgerManager := InitializeLedgerManager(client, pool)
	ledger.RegisterModule(client, ledgerManager)

	e621Client := e621.NewE621Service("twotto/1.0 (DownloadableFox)")
	e621.RegisterModule(client, e621Client)

	/*
		web := InitializeFiberServer()
		remote.RegisterModule(client, web)
	*/

	return nil
}
