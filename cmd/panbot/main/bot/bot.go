package bot

import (
	"github.com/brandao07/panbot/pkg/errors"
	"github.com/brandao07/panbot/pkg/todolist"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const botToken = "BOT_TOKEN"

var s *discordgo.Session
var storage *todolist.Storage

func init() {
	var err error
	errors.Check(nil, godotenv.Load())
	storage = todolist.NewStorage("items.json")
	s, err = discordgo.New("Bot " + os.Getenv(botToken))
	errors.Check(nil, err)
	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	addMessageHandler(s)
}

func Run() {
	defer func(s *discordgo.Session) {
		err := s.Close()
		errors.Check(nil, err)
	}(s)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := startSession(s)
	errors.Check(nil, err)

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	waitForShutdown()
	log.Println("Gracefully shutting down...")
}

func startSession(s *discordgo.Session) error {
	err := s.Open()
	if err != nil {
		return err
	}
	return nil
}

func waitForShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
