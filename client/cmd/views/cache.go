package views

import (
	"fmt"
	"net/rpc"

	"github.com/Malwarize/retro/client/controller"
	"github.com/Malwarize/retro/shared"
)

func CacheDisplay(client *rpc.Client) {
	songs := controller.GetCachedMusics(client)
	if len(songs) == 0 {
		fmt.Println("No music in cache")
		return
	}
	fmt.Println(GetTheme().ProgressStyle.Render("📁 Cache :"))
	for _, song := range songs {
		fmt.Printf(
			"\n    %s : %s\n", song.Hash[:shared.HashPrefixLength], song.Name,
		)
	}
}
