package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cristalhq/aconfig"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	BotConfig Config
)

func init() {
	// Load config
	loader := aconfig.LoaderFor(&BotConfig, aconfig.Config{})
	if err := loader.Load(); err != nil {
		panic(err)
	}

	// Set log level
	if BotConfig.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	client, err := discordgo.New("Bot " + BotConfig.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Discord client!")
	}

	// Bootstrap
	if err := bootstrap(client, &BotConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to bootstrap bot!")
	}

	// Run til end
	if err := client.Open(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Discord!")
	}
	defer client.Close()

	log.Info().Msg("[Main] Bot is set and running!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Warn().Msg("[Main] Stop signal sent! Stopping bot now...")

	// Remove all commands
	core.UnregisterAllCommands(client)
}
