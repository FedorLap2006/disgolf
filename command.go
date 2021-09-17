package disgolf

import (
	"github.com/bwmarrin/discordgo"
)

// A Handler processes the command
type Handler interface {
	HandleCommand(ctx *Ctx)
}

// HandlerFunc is a wrapper around Handler for functions
type HandlerFunc func(ctx *Ctx)

// HandleCommand implements Handler interface and calls the function with provided context
func (f HandlerFunc) HandleCommand(ctx *Ctx) {
	f(ctx)
}

// Command represents a command and extends discordgo.ApplicationCommand.
type Command struct {
	*discordgo.ApplicationCommand
	Handler Handler
}
