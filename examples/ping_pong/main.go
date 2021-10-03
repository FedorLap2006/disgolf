package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load .env: %w", err))
	}
}

func main() {
	bot, err := disgolf.New(os.Getenv("BOT_TOKEN"))
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
		MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
			_, _ = ctx.Reply("Hi, I'm a bot built on Disgolf library", true)
		}),

		Middlewares: []disgolf.Handler{
			disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
				fmt.Println("Middleware worked!")
				ctx.Next()
			}),
		},

		MessageMiddlewares: []disgolf.MessageHandler{
			disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
				fmt.Println("Message niddleware worked!", ctx.Arguments)
				ctx.Next()
			}),
		},
	})
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
		Prefixes:      []string{"d.", "dis.", "disgolf."},
		MentionPrefix: true,
	}))

	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("open exited with a error: %w", err))
	}
	defer bot.Close()
	err = bot.Router.Sync(bot.Session, "", "679281186975252480")
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
