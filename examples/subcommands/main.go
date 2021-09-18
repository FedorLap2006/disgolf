package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

var (
	token = flag.String("token", "", "Bot token")
)

func init() {
	flag.Parse()
}

func main() {
	bot, err := disgolf.New(*token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Router.Register(&disgolf.Command{
		Name:        "subcommands",
		Description: "Lo and behold, subcommands are coming!",
		Type:        discordgo.ChatApplicationCommand,
		SubCommands: disgolf.NewRouter([]*disgolf.Command{
			{
				Name:        "group",
				Description: "Subcommand group",
				SubCommands: disgolf.NewRouter([]*disgolf.Command{
					{
						Name:        "subcommand",
						Description: "Subcommand in a subcommand group",
						Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
							_ = ctx.Respond(&discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{Content: "hi (group)"},
							})
						}),
					},
				}),
			},
			{
				Name:        "subcommand",
				Description: "Just a subcommand",
				Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
					_ = ctx.Respond(&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "hi"},
					})
				}),
			},
		}),
	})
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	bot.AddHandler(bot.Router.HandleInteraction)

	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("open exited with a error: %w", err))
	}
	defer bot.Close()
	err = bot.Router.Sync(bot.Session, "", "GUILD-ID")
	if err != nil {
		log.Fatal(fmt.Errorf("cannot publish commands: %w", err))
	}

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
end:
	for {
		select {
		case <-stchan:
			break end
		default:
		}
		time.Sleep(time.Second)
	}
}
