package discord

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "start",
		Description: "Start Server",
	},
	{
		Name:        "send",
		Description: "Send command to server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "command",
				Description: "Command to execute",
				Required:    true,
			},
		},
	},
	{
		Name:        "ping",
		Description: "Check if server is up",
	},
}
