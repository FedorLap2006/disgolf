package disgolf

import (
	"github.com/bwmarrin/discordgo"
)

// OptionsMap is an alias for map of discordgo.ApplicationCommandInteractionDataOption
type OptionsMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

// Ctx is a context provided to a command. It embeds session for easier use,
// and contains interaction and preprocessed options.
type Ctx struct {
	*discordgo.Session
	Caller      *Command
	Interaction *discordgo.Interaction
	Options     OptionsMap
	OptionsRaw  []*discordgo.ApplicationCommandInteractionDataOption

	remainingHandlers []Handler
}

// Respond is a wrapper for ctx.Session.InteractionRespond
func (ctx *Ctx) Respond(response *discordgo.InteractionResponse) error {
	return ctx.Session.InteractionRespond(ctx.Interaction, response)
}

func (ctx *Ctx) Next() {
	if len(ctx.remainingHandlers) == 0 {
		return
	}

	handler := ctx.remainingHandlers[0]
	ctx.remainingHandlers = ctx.remainingHandlers[1:]

	handler.HandleCommand(ctx)
}

// NewCtx constructs ctx from given parameters.
func NewCtx(s *discordgo.Session, caller *Command, i *discordgo.Interaction, parent *discordgo.ApplicationCommandInteractionDataOption, handlers []Handler) *Ctx {
	options := i.ApplicationCommandData().Options
	if parent != nil {
		options = parent.Options
	}
	return &Ctx{
		Session:     s,
		Caller:      caller,
		Interaction: i,
		Options:     makeOptionMap(options),
		OptionsRaw:  options,

		remainingHandlers: handlers,
	}
}

func makeOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) (m OptionsMap) {
	m = make(OptionsMap, len(options))

	for _, option := range options {
		m[option.Name] = option
	}

	return
}

// MessageCtx is a context provided to message command handler. It contains the message
type MessageCtx struct {
	*discordgo.Session
	Caller    *Command
	Message   *discordgo.Message
	Arguments []string

	remainingHandlers []MessageHandler
}

// Next calls the next middleware / command handler.
func (ctx *MessageCtx) Next() {
	if len(ctx.remainingHandlers) == 0 {
		return
	}

	handler := ctx.remainingHandlers[0]
	ctx.remainingHandlers = ctx.remainingHandlers[1:]

	handler.HandleMessageCommand(ctx)
}

// Reply sends and returns a simple (content-only) message replying to the command message. If mention is true the command author is mentioned in the reply.
// It is a wrapper for discordgo.Session.ChannelMessageSendReply.
func (ctx *MessageCtx) Reply(content string, mention bool) (*discordgo.Message, error) {
	return ctx.ReplyComplex(&discordgo.MessageSend{
		Content: content,
	}, mention)
}

// ReplyComplex sends and returns a complex (with embds, attachments, etc) message replying to the command message.
// If mention is true the command author is mentioned in the reply. It is a wrapper for discordgo.Session.ChannelMessageSendComplex
func (ctx *MessageCtx) ReplyComplex(message *discordgo.MessageSend, mention bool) (*discordgo.Message, error) {
	message.Reference = ctx.Message.Reference()
	// TODO: https://github.com/bwmarrin/discordgo/pull/1009
	// message.AllowedMentions.RepliedUser = mention

	return ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, message)
}

// NewMessageCtx constructs context from a message.
// If argdelim is not empty it is a delimiter for the arguments, otherwise the arguments are split by a space.
func NewMessageCtx(s *discordgo.Session, caller *Command, m *discordgo.Message, arguments []string, handlers []MessageHandler) *MessageCtx {
	return &MessageCtx{
		Session:           s,
		Caller:            caller,
		Message:           m,
		Arguments:         arguments,
		remainingHandlers: handlers,
	}
}
