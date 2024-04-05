package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ConfigFile string
	BotConfig  *Config
)

func init() {
	// Parse flags
	flag.StringVar(&ConfigFile, "config", "config.json", "Path to the config file")
	flag.Parse()

	// Load config
	var err error
	if BotConfig, err = LoadConfig(ConfigFile); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config file!")
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
		panic(err)
	}

	// Bootstrap
	if err := bootstrap(client, BotConfig); err != nil {
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
	commands, err := client.ApplicationCommands(client.State.User.ID, "")
	if err != nil {
		log.Error().Err(err).Msg("[Main] Failed to list all commands!")
	}

	for _, command := range commands {
		if err := client.ApplicationCommandDelete(client.State.User.ID, "", command.ID); err != nil {
			log.Error().Err(err).Msgf("[Main] Failed to remove command %s!", command.Name)
		}
	}
}
