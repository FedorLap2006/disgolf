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
func (f HandlerFunc) HandleCommand(ctx *Ctx) { f(ctx) }

// A MessageHandler processes the message command
type MessageHandler interface {
	HandleMessageCommand(ctx *MessageCtx)
}

// HandlerFunc is a wrapper around MessageHandler for functions
type MessageHandlerFunc func(ctx *MessageCtx)

// HandleCommand implements MessageHandler interface and calls the function with provided context
func (f MessageHandlerFunc) HandleMessageCommand(ctx *MessageCtx) { f(ctx) }

// Command represents a command.
type Command struct {
	Name               string
	Description        string
	Options            []*discordgo.ApplicationCommandOption
	Type               discordgo.ApplicationCommandType
	Handler            Handler
	Middlewares        []Handler
	MessageHandler     MessageHandler
	MessageMiddlewares []MessageHandler

	// NOTE: nesting of more than 3 level has no effect
	SubCommands *Router
	// Custom payload for the command. Useful for module names, and such stuff.
	Custom interface{}
}

// ApplicationCommand converts Command to discordgo.ApplicationCommand.
func (cmd Command) ApplicationCommand() *discordgo.ApplicationCommand {
	applicationCommand := &discordgo.ApplicationCommand{
		Name:        cmd.Name,
		Description: cmd.Description,
		Options:     cmd.Options,
		Type:        cmd.Type,
	}
	for _, subcommand := range cmd.SubCommands.List() {
		applicationCommand.Options = append(applicationCommand.Options, subcommand.ApplicationCommandOption())
	}
	return applicationCommand
}

// ApplicationCommandOption converts Command to discordgo.ApplicationCommandOption (subcommand).
func (cmd Command) ApplicationCommandOption() *discordgo.ApplicationCommandOption {
	applicationCommand := cmd.ApplicationCommand()
	typ := discordgo.ApplicationCommandOptionSubCommand

	if cmd.SubCommands != nil && cmd.SubCommands.Count() != 0 {
		typ = discordgo.ApplicationCommandOptionSubCommandGroup
	}

	return &discordgo.ApplicationCommandOption{
		Name:        applicationCommand.Name,
		Description: applicationCommand.Description,
		Options:     applicationCommand.Options,
		Type:        typ,
	}
}
