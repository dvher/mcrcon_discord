package discord

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

var (
	BotToken       = os.Getenv("BOT_TOKEN")
	RemoveCommands = strings.ToLower(os.Getenv("REMOVE_COMMANDS")) == "true"
)

var Sess *discordgo.Session
var RegisteredCommands []*discordgo.ApplicationCommand

func init() {
	var err error
	Sess, err = discordgo.New("Bot " + BotToken)

	if err != nil {
		log.Fatalln("Invalid bot token")
	}
}

func init() {
	err := Sess.Open()

	if err != nil {
		log.Fatalln("Could not open discord session")
	}

	if RemoveCommands {
		removeCommands()
	}

	Sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	log.Println("Adding commands")
	for _, command := range commands {
		_, err := Sess.ApplicationCommandCreate(Sess.State.User.ID, "", command)

		if err != nil {
			log.Fatalf("Cannot create %v command: %v", command.Name, err)
		}

	}
}

func removeCommands() {
	registeredCommands, err := Sess.ApplicationCommands(Sess.State.User.ID, "")

	if err != nil {
		log.Fatalf("Cannot fetch commands: %v", err)
	}

	for _, cmd := range registeredCommands {
		err := Sess.ApplicationCommandDelete(Sess.State.User.ID, "", cmd.ID)

		if err != nil {
			log.Fatalf("Cannot remove %v command: %v", cmd.Name, err)
		}
	}
}
