package disgolf

import "github.com/bwmarrin/discordgo"

// OptionsMap is an alias for map of discordgo.ApplicationCommandInteractionDataOption
type OptionsMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

// Ctx is a context provided to a command. It embeds session for easier use,
// and contains interaction and preprocessed options.
type Ctx struct {
	*discordgo.Session
	Interaction *discordgo.Interaction
	Options     OptionsMap
	OptionsRaw  []*discordgo.ApplicationCommandInteractionDataOption
}

// Respond is a wrapper for ctx.Session.InteractionRespond
func (ctx *Ctx) Respond(response *discordgo.InteractionResponse) error {
	return ctx.Session.InteractionRespond(ctx.Interaction, response)
}

// NewCtx constructs ctx from given parameters.
func NewCtx(s *discordgo.Session, i *discordgo.Interaction, parent *discordgo.ApplicationCommandInteractionDataOption) *Ctx {
	options := i.ApplicationCommandData().Options
	if parent != nil {
		options = parent.Options
	}
	return &Ctx{
		Session:     s,
		Interaction: i,
		Options:     makeOptionMap(options),
		OptionsRaw:  options,
	}
}

func makeOptionMap(options []*discordgo.ApplicationCommandInteractionDataOption) (m OptionsMap) {
	m = make(OptionsMap, len(options))

	for _, option := range options {
		m[option.Name] = option
	}

	return
}
