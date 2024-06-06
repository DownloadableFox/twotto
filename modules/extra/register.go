package extra

import (
	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

func RegisterModule(client *discordgo.Session) {
	client.AddHandler(core.HandleEvent(HandleTwitterLinkEvent))
}
