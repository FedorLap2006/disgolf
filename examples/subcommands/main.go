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
		Name:        "subcommands",
		Description: "Lo and behold, subcommands are coming!",
		Type:        discordgo.ChatApplicationCommand,
		MessageMiddlewares: []disgolf.MessageHandler{
			disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
				fmt.Println("middleware")
				ctx.Next()
			}),
		},
		Middlewares: []disgolf.Handler{
			disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
				fmt.Println("middleware")
				ctx.Next()
			}),
		},
		SubCommands: disgolf.NewRouter([]*disgolf.Command{
			{
				Name:        "group",
				Description: "Subcommand group",
				MessageMiddlewares: []disgolf.MessageHandler{
					disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
						fmt.Println("group middleware")
						ctx.Next()
					}),
				},
				Middlewares: []disgolf.Handler{
					disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
						fmt.Println("group middleware")
						ctx.Next()
					}),
				},
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
						MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
							_, _ = ctx.Reply("hi (group)", false)
						}),
						MessageMiddlewares: []disgolf.MessageHandler{
							disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
								fmt.Println("individual middleware")
								ctx.Next()
							}),
						},
						Middlewares: []disgolf.Handler{
							disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
								fmt.Println("individual middleware")
								ctx.Next()
							}),
						},
					},
				}),
				MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
					_, _ = ctx.Reply("hi (group default)", false)
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
				MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
					_, _ = ctx.Reply("hi", false)
				}),
				MessageMiddlewares: []disgolf.MessageHandler{
					disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
						fmt.Println("individual middleware (2nd level)")
						ctx.Next()
					}),
				},
				Middlewares: []disgolf.Handler{
					disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
						fmt.Println("individual middleware (2nd level)")
						ctx.Next()
					}),
				},
			},
		}),
		MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
			_, _ = ctx.Reply("hi (default)", false)
		}),
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
