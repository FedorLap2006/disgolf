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

// Get returns a command by specified name.
func (r *Router) Get(name string) *Command {
	command, ok := r.Commands.Load(name)
	if ok {
		return command.(*Command)
	}

	return nil
}

// Update updates the command and does all behind-the-scenes work.
func (r *Router) Update(name string, newcmd *Command) (cmd *Command, err error) {
	rawCmd, ok := r.Commands.Load(name)

	if !ok {
		return nil, ErrCommandNotExists
	}

	r.Commands.Store(name, newcmd)
	return rawCmd.(*Command), nil
}

// Unregister removes a command from router
func (r *Router) Unregister(name string) (command *Command, existed bool) {
	var rawCommand interface{}
	rawCommand, existed = r.Commands.LoadAndDelete(name)

	if existed {
		command = rawCommand.(*Command)
	}

	return
}

// List returns all registered commands
func (r Router) List() (list []*Command) {
	r.Commands.Range(func(key, value interface{}) bool {
		list = append(list, value.(*Command))
		return true
	})
	return
}

// Count returns amount of commands stored
func (r Router) Count() (c int) {
	r.Commands.Range(func(_, _ interface{}) bool {
		c++
		return true
	})
	return
}

// Sync syncs all the commands with Discord.
func (r Router) Sync(s *discordgo.Session, application, guild string) error {
	if application == "" {
		if s.State.User == nil {
			panic("cannot determine application id")
		}
		application = s.State.User.ID
	}
	var commands []*discordgo.ApplicationCommand
	r.Commands.Range(func(_, cmd interface{}) bool {
		commands = append(commands, cmd.(*Command).applicationCommand())
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
