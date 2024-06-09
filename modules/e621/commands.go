package e621

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/downloadablefox/twotto/core"
	"github.com/rs/zerolog/log"
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
			AddIntegerOption(func(s *core.IntegerOptionBuilder) {
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
				Value:  fmt.Sprintf("%.2f MB", float64(post.Size)/1024.0/1024.0),
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

	// Create file request
	req, err := http.NewRequest(http.MethodGet, post.URL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file := &discordgo.File{
		Name:   fmt.Sprintf("post-%d.%s", post.ID, post.Ext),
		Reader: res.Body,
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
		Files:  []*discordgo.File{file},
	}); err != nil {
		return err
	}

	return nil
}

func publishThread(s *discordgo.Session, channelID, messageID, tags string, posts []*E621Post) error {
	// Assume
	success := true

	// Create a thread to send the posts
	thr, err := s.MessageThreadStartComplex(channelID, messageID, &discordgo.ThreadStart{
		Name:                fmt.Sprintf("Posts with tags `%s`", tags),
		AutoArchiveDuration: 60, // 1 hour
		Invitable:           true,
	})
	if err != nil {
		return err
	}

	// Send the posts
	for _, post := range posts {
		s.ChannelTyping(thr.ID)

		embed := GeneratePostEmbed(post)
		req, err := http.NewRequest(http.MethodGet, post.URL, nil)
		if err != nil {
			log.Warn().Err(err).Msgf("[E621YiffSearchCommand] Failed to create request for post #%d (source: %s)", post.ID, post.URL)
			success = false
			continue
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Warn().Err(err).Msgf("[E621YiffSearchCommand] Failed to download post #%d (source: %s)", post.ID, post.URL)
			success = false
			continue
		}
		defer res.Body.Close()

		file := &discordgo.File{
			Name:   fmt.Sprintf("post-%d.%s", post.ID, post.Ext),
			Reader: res.Body,
		}

		if _, err := s.ChannelMessageSendComplex(thr.ID, &discordgo.MessageSend{
			Embed: embed,
			Files: []*discordgo.File{file},
		}); err != nil {
			log.Warn().Err(err).Msgf("[E621YiffSearchCommand] Failed to send post #%d (source: %s)", post.ID, post.URL)
			success = false
			continue
		}
	}

	if !success {
		s.ChannelMessageSendComplex(thr.ID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Failed to send some posts!",
				Description: "There was an issue sending some posts. These posts were omitted from the thread!",
				Color:       core.ColorWarning,
			}},
		})
	}

	return nil
}

func handleSearch(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	data := e.ApplicationCommandData().Options[0]

	// Get E621 service from context
	svc, ok := ctx.Value(E621ServiceKey).(IE621Service)
	if !ok || svc == nil {
		return ErrE621ServiceNotFound
	}

	// Get tags
	tags, err := core.GetStringOption(data.Options, "tags")
	if err != nil {
		return err
	}

	limit := core.GetIntegerDefaultOption(data.Options, "limit", 10)
	page := core.GetIntegerDefaultOption(data.Options, "page", 1)

	if page <= 0 {
		return errors.New("`page` param should be > 0")
	}

	startTime := time.Now()

	// Send the looking for posts embed
	msg, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{{
			Title:       "Looking for posts...",
			Description: "Searching for posts (this may take a while) ...",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tags",
					Value:  fmt.Sprintf("`%s`", tags),
					Inline: true,
				},
				{
					Name:   "Elapsed Time",
					Value:  fmt.Sprintf("<t:%d:R>", startTime.Unix()),
					Inline: true,
				},
			},
			Color: core.ColorInfo,
		}},
	})
	if err != nil {
		return err
	}

	posts, err := svc.SearchPosts(tags, limit, page)
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		// Update interaction
		if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{{
				Title:       "No posts found!",
				Description: "No posts were found with the given tags.\n**Note:**Some files may be too large to send (25MB limit).",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Tags",
						Value:  fmt.Sprintf("`%s`", tags),
						Inline: true,
					},
					{
						Name:   "Elapsed Time",
						Value:  fmt.Sprintf("`%s`", time.Since(startTime).Round(time.Second).String()),
						Inline: true,
					},
				},
				Color: core.ColorError,
			}},
		}); err != nil {
			return err
		}

		return nil
	}

	// Send the posts to a thread
	if err := publishThread(s, msg.ChannelID, msg.ID, tags, posts); err != nil {
		// Operation cancelled
		s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{{
				Title:       "Operation cancelled! :(",
				Description: "There was an issue sending the posts to a thread.",
				Color:       core.ColorError,
			}},
		})

		return err
	}

	// Update interaction
	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{{
			Title:       "Posts sent!",
			Description: "Found posts have been sent below!\n**Note:**Some files might have been ommited as they were too large to send (25MB limit).",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Tags",
					Value:  fmt.Sprintf("`%s`", tags),
					Inline: false,
				},
				{
					Name:   "Count",
					Value:  fmt.Sprintf("%d", len(posts)),
					Inline: true,
				},
				{
					Name:   "Elapsed Time",
					Value:  fmt.Sprintf("%s sec", time.Since(startTime).Round(time.Second).String()),
					Inline: true,
				},
			},
			Color: core.ColorSuccess,
		}},
	}); err != nil {
		return err
	}

	return nil
}

func handlePost(ctx context.Context, s *discordgo.Session, e *discordgo.InteractionCreate) error {
	data := e.ApplicationCommandData().Options[0]

	// Get E621 service from context
	svc, ok := ctx.Value(E621ServiceKey).(IE621Service)
	if !ok || svc == nil {
		return ErrE621ServiceNotFound
	}

	// Get the post
	postID, err := core.GetIntegerOption(data.Options, "id")
	if err != nil {
		return err
	}

	post, err := svc.GetPostByID(postID)
	if err != nil {
		return err
	}

	// Send the post
	embed := GeneratePostEmbed(post)
	req, err := http.NewRequest(http.MethodGet, post.URL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file := &discordgo.File{
		Name:   fmt.Sprintf("post-%d.%s", post.ID, post.Ext),
		Reader: res.Body,
	}

	if _, err := s.InteractionResponseEdit(e.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
		Files:  []*discordgo.File{file},
	}); err != nil {
		return err
	}

	return nil
}
