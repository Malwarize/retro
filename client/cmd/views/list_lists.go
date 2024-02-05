package views

import (
	"fmt"
	"net/rpc"

	"github.com/Malwarize/goplay/client/controller"
)

func PlayListsDisplay(client *rpc.Client) {
	playlists := controller.PlayListsNames(client)
	if len(playlists) == 0 {
		fmt.Println("No playlists")
		return
	}
	fmt.Println(playListNameStyle.Render("📼 Playlists:	"))
	for _, playlist := range playlists {
		fmt.Printf("\n   - %s\n", playlist)
	}

}

func PlayListSongsDisplay(name string, client *rpc.Client) {
	songs := controller.PlayListSongs(name, client)
	if len(songs) == 0 {
		fmt.Println("No songs in playlist")
		return
	}
	fmt.Println(playListNameStyle.Render("🎧 Playlist: ") + name)
	for index, song := range songs {
		fmt.Printf(
			"\n    %d : %s\n", index, parseName(song),
		)
	}
}
