package main

import (
	"fmt"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

type ExampleModule struct{}

func (ExampleModule) PingFunctional(s *discordgo.Session) time.Duration {
	return s.HeartbeatLatency()
}

func (m ExampleModule) Ping(ctx *disgolf.Ctx) {
	_ = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				":ping_pong: %v",
				ctx.HeartbeatLatency(),
			),
		},
	})
}

func (ExampleModule) PingMessage(ctx *disgolf.MessageCtx) {
	_, _ = ctx.Reply(
		fmt.Sprintf(
			":ping_pong: %v",
			ctx.HeartbeatLatency(),
		),
		false,
	)
}

func (m ExampleModule) Commands() []*disgolf.Command {
	return []*disgolf.Command{
		{
			Name:           "ping",
			Description:    "Get bot ping",
			Handler:        disgolf.HandlerFunc(m.Ping),
			MessageHandler: disgolf.MessageHandlerFunc(m.PingMessage),
		},
		{
			Name:        "ping_functional",
			Description: "Get bot ping",
			Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
				_ = ctx.Respond(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							":ping_pong: %v",
							m.PingFunctional(ctx.Session),
						),
					},
				})
			}),
			MessageHandler: disgolf.MessageHandlerFunc(func(ctx *disgolf.MessageCtx) {
				_, _ = ctx.Reply(
					fmt.Sprintf(
						":ping_pong: %v",
						m.PingFunctional(ctx.Session),
					),
					false,
				)
			}),
		},
	}
}
