package views

import (
	"fmt"
	"github.com/Malwarize/retro/shared"
	"net/rpc"

	"github.com/Malwarize/retro/client/controller"
)

//	func CacheDisplay(client *rpc.Client) {
//		songs := controller.GetCachedMusics(client)
//		if len(songs) == 0 {
//			fmt.Println("No music in cache")
//			return
//		}
//		fmt.Println(GetTheme().ProgressStyle.Render("📁 Cache :"))
//		for _, song := range songs {
//			fmt.Printf(
//				"    %s : %s\n", song.Hash[:shared.HashPrefixLength], song.Name,
//			)
//		}
//	}
func CacheDisplay(client *rpc.Client) {
	songs := controller.GetCachedMusics(client)
	if len(songs) == 0 {
		fmt.Println("No music in cache")
		return
	}
	fmt.Print(GetTheme().PositionStyle.Render("📁 Cache\n"))
	fmt.Print("\n")

	for l, _ := range songs {
		if l == len(songs)-1 {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("└──["))
		} else {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("├──["))
		}
		fmt.Print(songs[l].Hash[:shared.HashPrefixLength])
		fmt.Print(GetTheme().ColoredTextStyle.Render("] "))
		fmt.Println(songs[l].Name)
	}

}
