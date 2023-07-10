package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// authenticate will try and authenticate to discord.
func authenticate(token string) (*discordgo.Session, error) {
	client, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	client.Identify.Intents = discordgo.IntentsGuildMessageTyping

	return client, nil
}

// authenticateWithBotToken authenticates you as a bot to Service via the given access token.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func authenticateWithBotToken(token string) (*discordgo.Session, error) {
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}

	return authenticate(token)
}

// authenticateWithOAuth2Token authenticates you to Service via the given OAUTH2 token.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func authenticateWithOAuth2Token(token string) (*discordgo.Session, error) {
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	return authenticate(token)
}
