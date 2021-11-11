package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

var modules = []interface{ Commands() []*disgolf.Command }{
	ExampleModule{},
}

func loadModules(bot *disgolf.Bot) {
	for _, m := range modules {
		for _, command := range m.Commands() {
			bot.Router.Register(command)
		}
	}
}

func main() {
	bot, err := disgolf.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialise session: %w", err))
	}
	loadModules(bot)
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
		Prefixes:      []string{"d.", "dis.", "disgolf."},
		MentionPrefix: true,
	}))
	bot.AddHandler(func(*discordgo.Session, *discordgo.Ready) { log.Println("Ready!") })
	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open session: %w", err))
	}
	err = bot.Router.Sync(bot.Session, "", os.Getenv("TEST_GUILD_ID"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to sync commands: %w", err))
	}

	ech := make(chan os.Signal)
	signal.Notify(ech, os.Kill, syscall.SIGTERM) //nolint:staticcheck
	<-ech
	_ = bot.Close()
}
