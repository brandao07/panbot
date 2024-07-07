package bot

import (
	"github.com/brandao07/panbot/pkg/todolist"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"sync"
	"time"
)

func addCommandHandler(s *discordgo.Session) {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create-item":            createItem,
		"delete-item":            deleteItem,
		"find-items-by-category": findItemsByCategory,
		"mark-as-completed":      markAsCompleted,
	}

	var wg sync.WaitGroup

	// Add Handlers
	for _, handler := range commandHandlers {
		wg.Add(1)
		go func(handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
			defer wg.Done()
			s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
					h(s, i)
				}
			})
		}(handler)
	}

	wg.Wait()
}

func addCommands(s *discordgo.Session) {
	categories := "(Anime, Book, Movie, Music Album, Song, TV Show)"

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-item",
			Description: "Creates a new item",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Item name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: categories,
					Required:    true,
				},
			},
		},
		{
			Name:        "delete-item",
			Description: "Deletes a item",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Item name",
					Required:    true,
				},
			},
		},
		{
			Name:        "find-items-by-category",
			Description: "Finds items by given category",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: categories,
					Required:    true,
				},
			},
		},
		{
			Name:        "mark-as-completed",
			Description: "Mark an item as completed",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Item name",
					Required:    true,
				},
			},
		},
	}
	// Adding commands
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func markAsCompleted(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name, ok := optionMap["name"]
	if !ok {
		replyWithError(s, i, "No name specified", "")
		return
	}

	item, err := storage.FindByName(name.StringValue())
	if err != nil {
		replyWithError(s, i, "Error finding the item", err.Error())
		return
	}

	err = storage.MarkAsCompleted(item)
	if err != nil {
		replyWithError(s, i, "Error completing the item", err.Error())
		return
	}
	reply(s, i, "Item completed successfully", nil)
}

func findItemsByCategory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	category, ok := optionMap["category"]
	if !ok {
		replyWithError(s, i, "No category specified", "")
		return
	}

	items, err := storage.FindByCategory(category.StringValue())
	if err != nil {
		replyWithError(s, i, "Error finding items by category", err.Error())
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name: "Name",
			Value: func(items *[]todolist.ItemDTO) string {
				var names []string
				for _, item := range *items {
					names = append(names, item.Name)
				}
				return strings.Join(names, "\n")
			}(items),
			Inline: true,
		},
		{
			Name: "Added by",
			Value: func(items *[]todolist.ItemDTO) string {
				var users []string
				for _, item := range *items {
					users = append(users, item.User)
				}
				return strings.Join(users, "\n")
			}(items),
			Inline: true,
		},
	}

	reply(s, i, category.StringValue()+" items", fields)
}

func deleteItem(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name, ok := optionMap["name"]
	if !ok {
		replyWithError(s, i, "No name specified", "")
		return
	}

	item, err := storage.FindByName(name.StringValue())
	if err != nil {
		replyWithError(s, i, "Error finding the item", err.Error())
		return
	}

	err = storage.Remove(item)
	if err != nil {
		replyWithError(s, i, "Error removing the item", err.Error())
		return
	}
	reply(s, i, "Item deleted successfully", nil)
}

func createItem(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	username := i.Member.User.Username

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name, ok := optionMap["name"]
	if !ok {
		replyWithError(s, i, "No name specified", "")
		return
	}

	category, ok := optionMap["category"]
	if !ok {
		replyWithError(s, i, "No category specified", "")
		return
	}

	item, err := todolist.NewItem(name.StringValue(), category.StringValue(), username)
	if err != nil {
		replyWithError(s, i, "Error creating the item", err.Error())
		return
	}

	err = storage.Add(item)
	if err != nil {
		replyWithError(s, i, "Error saving the item", err.Error())
		return
	}

	reply(s, i, "Item added successfully", nil)
}

func replyWithError(s *discordgo.Session, i *discordgo.InteractionCreate, title, description string) {

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xff0000,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: description,
		Title:       title,
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	})
}

func reply(s *discordgo.Session, i *discordgo.InteractionCreate, content string, fields []*discordgo.MessageEmbedField) {
	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0x00ff00,
		Fields:    fields,
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     content,
	}
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	})
}
