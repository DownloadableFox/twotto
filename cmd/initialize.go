//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/modules/debug"
	"github.com/downloadablefox/twotto/modules/extra"
	"github.com/downloadablefox/twotto/modules/ledger"
	"github.com/downloadablefox/twotto/modules/whitelist"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Token       string `json:"token"`
	Debug       bool   `json:"debug"`
	DatabaseURL string `json:"database_url"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err = json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
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

func InitializeWhitelistManager(pool *pgxpool.Pool) whitelist.WhitelistManager {
	wire.Build(whitelist.NewPostgresWhitelistManager)
	return nil
}

func InitializeLedgerManager(session *discordgo.Session, pool *pgxpool.Pool) ledger.LedgerManager {
	wire.Build(ledger.NewRepoLedgerManager, ledger.NewLedgerPostgresRepository)
	return nil
}

func bootstrap(client *discordgo.Session, config *Config) error {
	// Set intents
	client.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentGuildMessageReactions | discordgo.IntentGuildMembers | discordgo.IntentGuildBans

	pool, err := InitializeDatabasePool(config)
	if err != nil {
		return err
	}

	// Register modules
	debug.RegisterModule(client)
	extra.RegisterModule(client)

	whitelistManager := InitializeWhitelistManager(pool)
	whitelist.RegisterModule(client, whitelistManager)

	ledgerManager := InitializeLedgerManager(client, pool)
	ledger.RegisterModule(client, ledgerManager)

	return nil
}
