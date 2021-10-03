package disgolf

import (
	"strings"

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
func NewCtx(s *discordgo.Session, caller *Command, i *discordgo.Interaction, parent *discordgo.ApplicationCommandInteractionDataOption) *Ctx {
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

		remainingHandlers: append(caller.Middlewares, caller.Handler),
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
	Message   *discordgo.Message
	Arguments []string
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
func NewMessageCtx(s *discordgo.Session, m *discordgo.Message, argdelim string) *MessageCtx {
	if argdelim == "" {
		argdelim = " "
	}

	return &MessageCtx{
		Session:   s,
		Message:   m,
		Arguments: strings.Split(m.Content, argdelim),
	}
}
