package bot

import (
	"github.com/brandao07/panbot/pkg/errors"
	"github.com/bwmarrin/discordgo"
)

type messageType string

const (
	messagePing  messageType = "!ping"
	messageHello messageType = "!hello"
	messageBye   messageType = "!bye"
)

func addMessageHandler(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		switch messageType(m.Content) {
		case messagePing:
			sendMessage(s, m.ChannelID, "pong")
		case messageHello:
			sendMessage(s, m.ChannelID, "hi there!")
		case messageBye:
			sendMessage(s, m.ChannelID, "goodbye!")
		}
	})
}

func sendMessage(s *discordgo.Session, channelID string, content string) {
	_, err := s.ChannelMessageSend(channelID, content)
	errors.Check(nil, err)
}
