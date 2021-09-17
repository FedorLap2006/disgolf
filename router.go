package disgolf

import (
	"github.com/bwmarrin/discordgo"
)

// A Router stores all the commands and routes the interactions
type Router struct {
	// Commands is a map of registered commands.
	// Key is command name. Value is command instance.
	//
	// NOTE: it is not recommended to use it directly, use Register, Get, Update, Unregister functions instead.
	Commands map[string]*Command
}

// Register registers the command.
func (r *Router) Register(cmd *Command) {
	if _, ok := r.Commands[cmd.Name]; !ok {
		r.Commands[cmd.Name] = cmd
	}
}

// Get returns a command by specified name.
func (r *Router) Get(name string) *Command {
	return r.Commands[name]
}

// Update updates the command and does all behind-the-scenes work.
func (r *Router) Update(name string, newcmd *Command) (cmd *Command, err error) {

	if cmd, ok := r.Commands[name]; ok {
		r.Commands[name] = newcmd
		return cmd, nil
	}

	return nil, ErrCommandNotExists
}

// Unregister removes a command from router
func (r *Router) Unregister(name string) (command *Command, existed bool) {
	command, existed = r.Commands[name]

	if existed {
		delete(r.Commands, name)
	}

	return
}

// List returns all registered commands
func (r Router) List() (list []*Command) {
	for _, c := range r.Commands {
		list = append(list, c)
	}
	return
}

// Count returns amount of commands stored
func (r Router) Count() (c int) {
	return len(r.Commands)
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
	for _, c := range r.Commands {
		commands = append(commands, c.ApplicationCommand)

	}
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

	rcmd := r.Get(data.Name)
	if rcmd == nil {
		return
	}

	rcmd.Handler.HandleCommand(NewCtx(s, i.Interaction, nil))
}

// NewRouter constructs a router from a set of predefined commands.
func NewRouter(initial []*Command) (r *Router) {
	r = new(Router)
	for _, cmd := range initial {
		r.Register(cmd)
	}

	return
}
