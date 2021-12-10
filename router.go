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

func (r *Router) getSubcommand(cmd *Command, opt *discordgo.ApplicationCommandInteractionDataOption, parent []Handler) (*Command, *discordgo.ApplicationCommandInteractionDataOption, []Handler) {
	if cmd == nil {
		return nil, nil, nil
	}

	subcommand := cmd.SubCommands.Get(opt.Name)
	switch opt.Type {
	case discordgo.ApplicationCommandOptionSubCommand:
		return subcommand, opt, append(parent, append(subcommand.Middlewares, subcommand.Handler)...)
	case discordgo.ApplicationCommandOptionSubCommandGroup:
		return r.getSubcommand(subcommand, opt.Options[0], append(parent, subcommand.Middlewares...))
	}

	return cmd, nil, append(parent, cmd.Handler)
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
	handlers := append(cmd.Middlewares, cmd.Handler)
	if len(data.Options) != 0 {
		cmd, parent, handlers = r.getSubcommand(cmd, data.Options[0], cmd.Middlewares)
	}

	if cmd != nil {
		ctx := NewCtx(s, cmd, i.Interaction, parent, handlers)
		ctx.Next()
	}
}

type MessageHandlerConfig struct {
	// Prefixes got will respond to
	Prefixes      []string
	MentionPrefix bool

	ArgumentDelimiter string
}

func (r *Router) getMessageSubcommand(cmd *Command, arguments []string, parent []MessageHandler) (*Command, []string, []MessageHandler) {
	if len(arguments) == 0 {
		return cmd, arguments, append(parent, cmd.MessageHandler)
	}
	subcommand := cmd.SubCommands.Get(arguments[0])
	if subcommand != nil {
		if len(arguments) > 1 {
			return r.getMessageSubcommand(subcommand, arguments[1:], append(parent, subcommand.MessageMiddlewares...)) // TODO: opt-out
		} else {
			return subcommand, arguments[1:], append(parent, append(subcommand.MessageMiddlewares, subcommand.MessageHandler)...)
		}
	}
	return cmd, arguments, append(parent, cmd.MessageHandler)
}

func (r *Router) MakeMessageHandler(cfg *MessageHandlerConfig) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	if cfg.ArgumentDelimiter == "" {
		cfg.ArgumentDelimiter = " "
	}
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

		arguments := strings.Split(m.Content, cfg.ArgumentDelimiter)

		commandName := arguments[0]

		command, ok := r.Commands[commandName]
		if !ok {
			return
		}
		arguments = arguments[1:]

		command, arguments, handlers := r.getMessageSubcommand(command, arguments, command.MessageMiddlewares)
		if command.MessageHandler == nil {
			return
		}

		ctx := NewMessageCtx(s, command, m.Message, arguments, handlers)
		ctx.Next()
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
