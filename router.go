package disgolf

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// A Router stores all the commands and routes the interactions
type Router struct {
	// Commands is a map of registered commands.
	// Key is command name. Value is command instance.
	//
	// NOTE: it is not recommended to use it directly, use Register, Get, Update, Unregister functions instead.
	Commands map[string]*Command

	Syncer CommandSyncer
}

// Register registers the command.
func (r *Router) Register(cmd *Command) {
	if _, ok := r.Commands[cmd.Name]; !ok {
		r.Commands[cmd.Name] = cmd
	}
}

// Get returns a command by specified name.
func (r *Router) Get(name string) *Command {
	if r == nil {
		return nil
	}
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
func (r *Router) List() (list []*Command) {
	if r == nil {
		return nil
	}

	for _, c := range r.Commands {
		list = append(list, c)
	}
	return
}

// Count returns amount of commands stored
func (r *Router) Count() (c int) {
	if r == nil {
		return 0
	}
	return len(r.Commands)
}

// A CommandSyncer syncs all the commands with Discord.
type CommandSyncer interface {
	Sync(r *Router, s *discordgo.Session, application, guild string) error
}

// BulkCommandSyncer syncs all the commands using ApplicationCommandBulkOverwrite function.
type BulkCommandSyncer struct{}

// Sync implements CommandSyncer interface.
func (BulkCommandSyncer) Sync(r *Router, s *discordgo.Session, application string, guild string) error {
	if application == "" {
		panic("empty application id")
	}

	var commands []*discordgo.ApplicationCommand
	for _, c := range r.Commands {
		commands = append(commands, c.ApplicationCommand())
	}
	_, err := s.ApplicationCommandBulkOverwrite(application, guild, commands)
	return err
}

// Sync wraps Router.Syncer and automatically detects application id.
func (r *Router) Sync(s *discordgo.Session, application, guild string) error {
	if application == "" {
		if s.State.User == nil {
			panic("cannot determine application id")
		}
		application = s.State.User.ID
	}
	return r.Syncer.Sync(r, s, application, guild)
}

func (r *Router) getSubcommand(cmd *Command, opt *discordgo.ApplicationCommandInteractionDataOption) (*Command, *discordgo.ApplicationCommandInteractionDataOption) {
	if cmd == nil {
		return nil, nil
	}

	switch opt.Type {
	case discordgo.ApplicationCommandOptionSubCommand:
		return cmd.SubCommands.Get(opt.Name), opt
	case discordgo.ApplicationCommandOptionSubCommandGroup:
		return r.getSubcommand(cmd.SubCommands.Get(opt.Name), opt.Options[0])
	}

	return cmd, nil
}

// HandleInteraction is an interaction handler passed to discordgo.Session.AddHandler.
func (r *Router) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()

	cmd := r.Get(data.Name)
	if cmd == nil {
		return
	}
	var parent *discordgo.ApplicationCommandInteractionDataOption
	if len(data.Options) != 0 {
		cmd, parent = r.getSubcommand(cmd, data.Options[0])
	}

	if cmd != nil {
		cmd.Handler.HandleCommand(NewCtx(s, i.Interaction, parent))
	}
}

type MessageHandlerConfig struct {
	// Prefixes got will respond to
	Prefixes      []string
	MentionPrefix bool

	ArgumentDelimiter string
}

func (r *Router) MakeMessageHandler(cfg *MessageHandlerConfig) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var match bool
		var prefixes []string
		prefixes = cfg.Prefixes
		if cfg.MentionPrefix {
			prefixes = append(prefixes,
				"<@"+s.State.User.ID+">",
				"<@!"+s.State.User.ID+">",
				"<@"+s.State.User.ID+"> ",
				"<@!"+s.State.User.ID+"> ",
			)
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(m.Content, prefix) {
				match = true
				m.Content = strings.TrimSpace(strings.TrimPrefix(m.Content, prefix))
				break
			}
		}

		if !match {
			return
		}
		ctx := NewMessageCtx(s, m.Message, cfg.ArgumentDelimiter)

		command := ctx.Arguments[0]

		handler, ok := r.Commands[command]

		if !ok || handler.MessageHandler == nil {
			return
		}

		ctx.Arguments = ctx.Arguments[1:]
		handler.MessageHandler.HandleMessageCommand(ctx)
	}
}

// NewRouter constructs a router from a set of predefined commands.
func NewRouter(initial []*Command) (r *Router) {
	r = &Router{Commands: make(map[string]*Command, len(initial)), Syncer: BulkCommandSyncer{}}
	for _, cmd := range initial {
		r.Register(cmd)
	}

	return
}
