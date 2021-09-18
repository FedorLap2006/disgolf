package main

import (
	"flag"
	"fmt"
	"log"

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
		Name:        "ping_pong",
		Description: "Ping it!",
		Type:        discordgo.ChatApplicationCommand,
		Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
			_ = ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hi, I'm a bot built on Disgolf library.",
				},
			})
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
	err = bot.Router.Sync(bot.Session, "", "TEST-GUILD-ID")
	if err != nil {
		log.Fatal(fmt.Errorf("cannot publish commands: %w", err))
	}

}
