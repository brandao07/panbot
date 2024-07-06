package bot

import (
	"github.com/brandao07/panbot/pkg/todolist"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
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
		respondToInteraction(s, i, "❌ No name specified")
		return
	}

	item, err := storage.FindByName(name.StringValue())
	if err != nil {
		respondToInteraction(s, i, "❌ Error finding the item:  "+err.Error())
		return
	}

	err = storage.MarkAsCompleted(item)
	if err != nil {
		respondToInteraction(s, i, "❌ Error completing the item: "+err.Error())
		return
	}
	respondToInteraction(s, i, "Item completed successfully ✅")
}

func findItemsByCategory(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO Implement this method
	respondToInteraction(s, i, "❌ Method not implemented!")
}

func deleteItem(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	name, ok := optionMap["name"]
	if !ok {
		respondToInteraction(s, i, "❌ No name specified")
		return
	}

	item, err := storage.FindByName(name.StringValue())
	if err != nil {
		respondToInteraction(s, i, "❌ Error finding the item: "+err.Error())
		return
	}

	err = storage.Remove(item)
	if err != nil {
		respondToInteraction(s, i, "❌ Error removing the item: "+err.Error())
		return
	}
	respondToInteraction(s, i, "Item deleted successfully 🗑️")
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
		respondToInteraction(s, i, "❌ No name specified")
		return
	}

	category, ok := optionMap["category"]
	if !ok {
		respondToInteraction(s, i, "❌ No category specified")
		return
	}

	item, err := todolist.NewItem(name.StringValue(), category.StringValue(), username)
	if err != nil {
		respondToInteraction(s, i, "❌ Error creating the item: "+err.Error())
		return
	}

	err = storage.Add(item)
	if err != nil {
		respondToInteraction(s, i, "❌ Error saving the item: "+err.Error())
		return
	}

	respondToInteraction(s, i, "Item added successfully ✨")
}

func respondToInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
