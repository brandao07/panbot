package main

import (
	"github.com/brandao07/panbot/cmd/panbot/main/bot"
	"github.com/brandao07/panbot/pkg/errors"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"os"
)

const botToken = "BOT_TOKEN"

func init() {
	color.Set(color.FgHiBlue)
	errors.Check(nil, godotenv.Load())
}

func main() {
	bot.Run(os.Getenv(botToken))
}
