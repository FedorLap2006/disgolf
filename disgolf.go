package disgolf

import (
	"github.com/bwmarrin/discordgo"
)

// A Bot wraps discordgo.Session with configuration and a router.
type Bot struct {
	*discordgo.Session

	Router *Router
}

// New constructs a Bot, from a authentication token.
func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session: session,
		Router:  NewRouter(nil),
	}, nil
}
