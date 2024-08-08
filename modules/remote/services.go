package remote

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	FiberServerKey   = core.NewIdentifier("remote", "service/fiber")
	ErrFiberNotFound = errors.New("fiber service not found in context (missing injection)")
)

func WebMidwareLogger(c *fiber.Ctx) error {
	log.Info().Msgf("[RemoteService] %s %s", c.Method(), c.Path())
	return c.Next()
}

func HandleHeartbeat(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

var (
	ActivityGuildID  = "1024188032829628497"
	ActivityMemberID = "556132236697665547"
)

func HandleGetActivity(s *discordgo.Session) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get presence of the user in the guild
		presence, err := s.State.Presence(ActivityGuildID, ActivityMemberID)
		if err != nil {
			log.Error().Err(err).Msg("[RemoteActivity]Failed to get presence from user!")
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to get presence from user!")
		}

		// Return the presence
		return c.JSON(fiber.Map{
			"status": "ok",
			"presence": fiber.Map{
				"status":     presence.Status,
				"activities": presence.Activities,
			},
		})
	}
}
