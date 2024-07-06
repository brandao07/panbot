package bot

import (
	"github.com/brandao07/panbot/pkg/todolist"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const botToken = "BOT_TOKEN"

var storage *todolist.Storage

func Run() {
	var wg sync.WaitGroup

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Create Storage
	storage = todolist.NewStorage("items.json")

	// Create Discord Bot
	s, err := discordgo.New("Bot " + os.Getenv(botToken))
	if err != nil {
		log.Fatal(err)
	}
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	wg.Add(1)

	go func(s *discordgo.Session) {
		// Discord Bot Handlers
		addCommandHandler(s)
		addMessageHandler(s)
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		})

		wg.Done()
	}(s)

	wg.Wait()

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s)

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Add commands
	addCommands(s)

	waitForShutdown()
	log.Println("Gracefully shutting down...")
}

func waitForShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
