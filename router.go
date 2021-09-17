package disgolf

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// A Router stores all the commands and routes the interactions
type Router struct {
	// Commands is a map of registered commands.
	// Key is a string - command name. Value is a pointer to a Command.
	//
	// NOTE: it is not recommended to use it directly, use Register, Get, Update, Unregister functions instead.
	Commands sync.Map
}

// Register registers the command.
func (r *Router) Register(cmd *Command) {
	r.Commands.Store(cmd.Name, cmd)
}

// Sync syncs all the commands with Discord.
func (r *Router) Sync(s *discordgo.Session, application, guild string) error {
	if application == "" {
		if s.State.User == nil {
			panic("cannot determine application id")
		}
		application = s.State.User.ID
	}
	var commands []*discordgo.ApplicationCommand
	r.Commands.Range(func(_, rawCmd interface{}) bool {
		commands = append(commands, rawCmd.(*Command).ApplicationCommand)
		return true
	})
	_, err := s.ApplicationCommandBulkOverwrite(application, guild, commands) // TODO: syncer
	return err
}

// HandleInteraction is an interaction handler passed to discordgo.Session.AddHandler.
//
func (r *Router) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()

	rcmd, ok := r.Commands.Load(data.Name)
	if !ok {
		return
	}

	cmd := rcmd.(*Command)
	cmd.Handler.HandleCommand(NewCtx(s, i.Interaction, nil))
}

// NewRouter constructs a router from a set of predefined commands.
func NewRouter(initial []*Command) (r *Router) {
	r = new(Router)
	for _, cmd := range initial {
		r.Register(cmd)
	}

	return
}
