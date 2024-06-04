package views

import (
	"fmt"
	"net/rpc"

	"github.com/Malwarize/retro/client/controller"
)

func PlayListsDisplay(client *rpc.Client) {
	playlists := controller.GetPlayListsNames(client)
	if len(playlists) == 0 {
		fmt.Println("No playlists")
		return
	}
	fmt.Println(GetTheme().PositionStyle.Render("📼 Playlists"))
	fmt.Println()

	for index, playlist := range playlists {
		if index == len(playlists)-1 {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("└──[ "))
		} else {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("├──[ "))
		}
		fmt.Println(playlist)
	}
}

func PlayListMusicsDisplay(name string, client *rpc.Client) {
	songs := controller.PlayListMusics(name, client)
	if len(songs) == 0 {
		fmt.Println("No songs in playlist")
		return
	}
	fmt.Println(GetTheme().PositionStyle.Render("🎧 Playlist: ") + name)
	fmt.Println()
	for index, song := range songs {
		if index == len(songs)-1 {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("└──[ "))
		} else {
			fmt.Print(GetTheme().PositionStyle.Copy().Inherit(GetTheme().ColoredTextStyle).Render("├──[ "))
		}
		fmt.Println(song)
	}
}
