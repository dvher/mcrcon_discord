package discord

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dvher/mcconbot/pkg/rcon"
	_ "github.com/joho/godotenv/autoload"
)

var handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"start": handleStart,
	"ping":  handlePing,
	"send":  handleSendCommand,
}

func handleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := os.Getenv("START_COMMAND")

	if command == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There's no defined command for starting the server!",
			},
		})
		return
	}

	commandArgs := strings.Split(command, " ")

	if len(commandArgs) == 1 {
		commandArgs = []string{commandArgs[0], ""}
	}

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	err := cmd.Run()

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Cannot run command %v", command),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Starting server",
		},
	})
}

func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conn, err := rcon.NewConnection()

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Got no response from server with error %s", err),
			},
		})
		return
	}

	defer conn.Close()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Server is up!",
		},
	})
}

func handleSendCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conn, err := rcon.NewConnection()

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not connect to server",
			},
		})
		return
	}

	defer conn.Close()

	command := i.ApplicationCommandData().GetOption("command")

	response, err := conn.SendCommand(command.StringValue())

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not send command to server",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: string(response.Payload),
		},
	})
}
