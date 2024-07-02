package bot

import (
	"github.com/brandao07/panbot/pkg/errors"
	"github.com/brandao07/panbot/pkg/todolist"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var storage *todolist.Storage

func Run(token string) {
	storage = todolist.NewStorage("items.json")
	sess, err := initSession(token)
	errors.Check(nil, err)

	defer func(sess *discordgo.Session) {
		err := sess.Close()
		errors.Check(nil, err)
	}(sess)

	addMessageHandler(sess)

	err = startSession(sess)
	errors.Check(nil, err)

	log.Println("Bot is now running.")

	waitForShutdown()
}

func initSession(token string) (*discordgo.Session, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	return sess, nil
}

func startSession(sess *discordgo.Session) error {
	err := sess.Open()
	if err != nil {
		return err
	}
	return nil
}

func waitForShutdown() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	color.Unset()
	<-sc
}
