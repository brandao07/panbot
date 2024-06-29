package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
)

var Token string

func Run() {
	// create a session
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic("error creating Discord session: " + err.Error())
	}

	discord.Open()
	defer discord.Close()

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
