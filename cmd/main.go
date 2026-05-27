package main

import (
	rcon "github.com/dvher/mcconbot/pkg"
)

func main() {
	conn := rcon.NewConnection("127.0.0.1", 25575, "Dvher2510%.")

	conn.SendCommand("list")

	defer conn.Close()
}
