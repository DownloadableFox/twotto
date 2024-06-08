package e621

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
)

// Now using the new command builder
var YiffCommand = core.NewCommandBuilder().
	SetName("yiff").
	SetDescription("Get a random yiff image from e621").
	SetDMPermission(true).
	SetNSFW(true).
	AddSubCommand(func(subcommand *core.SubCommandBuilder) {
		subcommand.SetName("random").
			SetDescription("Get a random furry image from e621")
	}).
	AddSubCommand(func(subcommand *core.SubCommandBuilder) {
		subcommand.SetName("search").
			SetDescription("Search for a furry image on e621").
			AddStringOption(func(s *core.StringOptionBuilder) {
				s.SetName("tags").
					SetDescription("The tags to search for space separated").
					SetRequired(true)
			}).
			AddIntegerOption(func(i *core.IntegerOptionBuilder) {
				i.SetName("limit").
					SetDescription("The maximum number of images to return").
					SetRequired(false)
			}).
			AddIntegerOption(func(i *core.IntegerOptionBuilder) {
				i.SetName("page").
					SetDescription("The page to return").
					SetRequired(false)
			})
	}).
	AddSubCommand(func(subcommand *core.SubCommandBuilder) {
		subcommand.SetName("post").
			SetDescription("Given an ID it retrives a post.").
			AddStringOption(func(s *core.StringOptionBuilder) {
				s.SetName("id").
					SetDescription("The ID of the post").
					SetRequired(true)
			})
	}).
	Build()

func GeneratePostEmbed(post *E621Post) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("E621 Post #%d", post.ID),
		Description: "You can find this post by clicking on the following URL.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Extension",
				Value:  fmt.Sprintf("`%s`", post.Ext),
				Inline: true,
			},
			{
				Name:   "Size",
				Value:  fmt.Sprintf("%.2f MB", float64(post.Size)/10000.0),
				Inline: true,
			},
		},
		URL:   fmt.Sprintf("https://e621.net/posts/%d", post.ID),
		Color: core.ColorInfo,
	}

	return embed
}

func HandleYiffCommand(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	data := e.ApplicationCommandData()

	// Defer the response
	if err := s.InteractionRespond(e.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	// Handle the subcommands
	switch core.GetSubcommandOption(data.Options) {
	case "random":
		return handleRandom(ctx, s, e)
	case "search":
		return handleSearch(ctx, s, e)
	case "post":
		return handlePost(ctx, s, e)
	}

	return nil
}

func handleRandom(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	// Get E621 service from context
	svc, ok := ctx.Value(E621ServiceKey).(IE621Service)
	if !ok || svc == nil {
		return ErrE621ServiceNotFound
	}

	// Get the post
	post, err := svc.GetRandomPost()
	if err != nil {
		return err
	}

	// Send the post
	embed := GeneratePostEmbed(post)
	attachment := &discordgo.MessageAttachment{
		URL:      post.URL,
		Filename: fmt.Sprintf("post-%d.%s", post.ID, post.Ext),
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Content:     &post.URL,
		Embeds:      &[]*discordgo.MessageEmbed{embed},
		Attachments: &[]*discordgo.MessageAttachment{attachment},
	}); err != nil {
		return err
	}

	return nil
}

func handleSearch(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return errors.New("not implemented")
}

func handlePost(_ context.Context, _ *discordgo.Session, _ *discordgo.InteractionCreate) error {
	return errors.New("not implemented")
}
