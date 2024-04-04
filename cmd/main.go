package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/modules/debug"
	"github.com/downloadablefox/twotto/modules/extra"
	"github.com/downloadablefox/twotto/modules/whitelist"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	client, err := discordgo.New("Bot + TOKEN")
	if err != nil {
		panic(err)
	}

	// Set intents
	client.Identify.Intents = discordgo.IntentGuildMessages

	// Register modules
	debug.RegisterModule(client)
	extra.RegisterModule(client)
	whitelist.RegisterModule(client)

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
