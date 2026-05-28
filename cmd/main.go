package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dvher/mcconbot/pkg/discord"
)

func main() {
	defer discord.Sess.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl-C to exit")
	<-stop
	log.Println("Shutting down...")
}
