package core

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func UnregisterAllCommands(s *discordgo.Session) {
	workers := 5
	wg := &sync.WaitGroup{}
	commandChan := make(chan *discordgo.ApplicationCommand, 100)

	// Remove all commands before restarting
	servers, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		log.Error().Err(err).Msg("[Core] Failed to list all guilds!")
	}

	go func() {
		defer close(commandChan)

		globalCommands, err := s.ApplicationCommands(s.State.User.ID, "")
		if err != nil {
			log.Error().Err(err).Msg("[Core] Failed to list all global commands!")
		}

		for _, command := range globalCommands {
			commandChan <- command
		}

		for _, server := range servers {
			commands, err := s.ApplicationCommands(s.State.User.ID, server.ID)
			if err != nil {
				log.Error().Err(err).Msg("[Core] Failed to list all commands!")
			}

			for _, command := range commands {
				commandChan <- command
			}
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for command := range commandChan {
				if err := s.ApplicationCommandDelete(s.State.User.ID, command.GuildID, command.ID); err != nil {
					log.Error().Err(err).Msg("[Core] Failed to delete command!")
				}
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()
}
