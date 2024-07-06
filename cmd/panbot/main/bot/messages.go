package bot

import (
	"github.com/bwmarrin/discordgo"
)

type messageType string

const (
	messagePing messageType = "!ping"
)

func addMessageHandler(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		switch messageType(m.Content) {
		case messagePing:
			sendMessage(s, m.ChannelID, "pong")
		}
	})
}

func sendMessage(s *discordgo.Session, channelID string, content string) {
	_, _ = s.ChannelMessageSend(channelID, content)

}
