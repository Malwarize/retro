package views

import (
	"fmt"
	"net/rpc"

	"github.com/Malwarize/retro/client/controller"
)

func CacheDisplay(client *rpc.Client) {
	songs := controller.GetCachedMusics(client)
	if len(songs) == 0 {
		fmt.Println("No music in cache")
		return
	}
	fmt.Println(GetTheme().ProgressStyle.Render("📁 Cache :"))
	for index, song := range songs {
		fmt.Printf(
			"\n    %d : %s\n", index, song,
		)
	}
}
