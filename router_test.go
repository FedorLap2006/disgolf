package disgolf_test

import (
	"testing"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

var router = disgolf.NewRouter(nil)

func TestRouter_Register(t *testing.T) {
	command := &disgolf.Command{
		Name:        "test_register",
		Description: "test register",
		Type:        discordgo.ChatApplicationCommand,
	}

	router.Register(command)
	defer router.Unregister(command.Name)

	if assert.NotNil(t, router.Commands[command.Name]) {
		assert.Equal(t, router.Commands[command.Name], command)
	}
}

func TestRouter_Get(t *testing.T) {
	command := &disgolf.Command{
		Name:        "test_get",
		Description: "test get",
		Type:        discordgo.ChatApplicationCommand,
	}

	router.Register(command)
	defer router.Unregister(command.Name)

	if assert.NotNil(t, router.Commands[command.Name]) {
		assert.Equal(t, router.Get(command.Name), command)
	}
}

func TestRouter_Update(t *testing.T) {
	command := &disgolf.Command{
		Name:        "test_update",
		Description: "test update",
		Type:        discordgo.MessageApplicationCommand,
	}
	router.Register(command)
	defer router.Unregister(command.Name)

	newCommand := &disgolf.Command{
		Name:        "test_update",
		Description: "test update",
		Type:        discordgo.ChatApplicationCommand,
	}

	oldcmd, err := router.Update("test_update", newCommand)

	assert.NoError(t, err)
	assert.Equal(t, command, oldcmd)
	assert.Equal(t, newCommand, router.Get(command.Name))
}
func TestRouter_Unregister(t *testing.T) {
	command := &disgolf.Command{
		Name:        "test_unregister",
		Description: "test unregister",
		Type:        discordgo.ChatApplicationCommand,
	}

	router.Register(command)
	defer router.Unregister(command.Name)

	if assert.NotNil(t, router.Get(command.Name)) {
		assert.Equal(t, command, router.Get(command.Name))
	}

	oldcmd, existed := router.Unregister(command.Name)
	assert.True(t, existed)
	assert.Equal(t, command, oldcmd)
}
func TestRouter_List(t *testing.T) {
	commandList := []*disgolf.Command{
		{
			Name:        "test_list_0",
			Description: "test list 0",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "test_list_1",
			Description: "test list 1",
			Type:        discordgo.MessageApplicationCommand,
		},
		{
			Name:        "test_list_2",
			Description: "test list 2",
			Type:        discordgo.UserApplicationCommand,
		},
	}

	for _, command := range commandList {
		router.Register(command)
		defer router.Unregister(command.Name)
	}

	assert.Equal(t, commandList, router.List())
	assert.Len(t, router.Commands, len(commandList))
}
func TestRouter_Count(t *testing.T) {
	commandList := []*disgolf.Command{
		{
			Name:        "test_count_0",
			Description: "test count 0",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "test_count_1",
			Description: "test count 1",
			Type:        discordgo.MessageApplicationCommand,
		},
		{
			Name:        "test_count_2",
			Description: "test count 2",
			Type:        discordgo.UserApplicationCommand,
		},
	}

	for _, command := range commandList {
		router.Register(command)
		defer router.Unregister(command.Name)
	}
	assert.Len(t, router.Commands, len(commandList))
	assert.Equal(t, len(commandList), router.Count())
}
