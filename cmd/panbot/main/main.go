package main

import (
	"github.com/brandao07/panbot/cmd/panbot/main/bot"
	"os"
)

func main() {
	bot.Token = os.Getenv("DISCORD_TOKEN")
	bot.Run()
}
