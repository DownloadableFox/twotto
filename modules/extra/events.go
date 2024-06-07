package extra

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
)

func HandleOnReadyEvent(ctx context.Context, s *discordgo.Session, e *discordgo.Ready) error {
	// Register commands
	err := core.ApplyCommands(
		SayCommand,
	).For(s, "")
	if err != nil {
		log.Warn().Err(err).Msg("[DebugModule] Failed to register commands")
		return err
	}

	return nil
}

var (
	TwitterLinkRegex = regexp.MustCompile(`(https?://)?(www\.)?(twitter|x)\.com/([a-zA-Z0-9_]+/status/[0-9]+)`)
)

// Event that catches a twitter/x link and fixes it with the correct content.
// replaces twitter.com or x.com with vxtwitter.com
func HandleTwitterLinkEvent(_ context.Context, s *discordgo.Session, e *discordgo.MessageCreate) error {
	// Ignore any bot messages
	if e.Author.ID == s.State.User.ID {
		return nil
	}

	// Check for substrings and store them in a strings list
	links := []string{}
	for _, match := range TwitterLinkRegex.FindAllStringSubmatch(e.Content, -1) {
		if len(match) != 5 {
			continue
		}

		links = append(links, match[4])
	}

	if len(links) == 0 {
		return nil
	}

	// Replace the links with the correct content
	linksStr := ""
	for _, link := range links {
		linksStr += fmt.Sprintf("https://vxtwitter.com/%s\n", link)
	}

	_, err := s.ChannelMessageSendReply(e.ChannelID, fmt.Sprintf("I have fixed the Twitter embeds for you ;3\n%s", linksStr), e.Reference())
	return err
}
