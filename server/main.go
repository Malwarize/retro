package main

import (
	"github.com/Malwarize/retro/config"
	"github.com/Malwarize/retro/server/player"
)

func main() {
	// load config
	cfg := config.GetConfig()

	player.StartIPCServer(cfg.ServerPort)
}
