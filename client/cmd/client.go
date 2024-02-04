package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Malwarize/goplay/client/cmd/views"
	"github.com/Malwarize/goplay/client/controller"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "play a song",
	Long:  `play a song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		if dir, err := cmd.Flags().GetString("dir"); err == nil && dir != "" {
			controller.PlayDir(dir, client)
			os.Exit(0)
		}

		if len(args) > 0 {
			song := strings.Join(args, " ")
			views.SearchThenSelect(song, client)
		} else {
			fmt.Println("no song specified")
			return
		}
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause the current song",
	Long:  `pause the current song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.Pause(client)
	},
}

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "resume the current song",
	Long:  `resume the current song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.Resume(client)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the current song",
	Long:  `stop the current song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.Stop(client)
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "play the next song", Long: `play the next song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.Next(client)
	},
}

var prevCmd = &cobra.Command{
	Use:   "prev",
	Short: "play the previous song",
	Long:  `play the previous song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.Prev(client)
	},
}
var seekCmd = &cobra.Command{Use: "seek",
	Short: "seek to a position in the current song",
	Long:  `seek to a position in the current song`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		var seekSeconds int
		if len(args) > 0 {
			var err error
			seekSeconds, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		controller.Seek(seekSeconds, client)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "get the current status of the player",
	Long:  `get the current status of the player`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		views.DisplayStatus(client)
	},
}

var playlistCmd = &cobra.Command{
	Use:   "list",
	Short: "list playlists",
	Long:  `list playlists`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			musics := controller.PlayListSongs(args[0], controller.GetClient())
			fmt.Println("Songs in playlist", args[0])
			for i, music := range musics {
				fmt.Printf("%d. %s\n", i, music)
			}
			return
		}
		client := controller.GetClient()
		lists := controller.PlayListsNames(client)
		fmt.Println("Playlists:")
		for i, list := range lists {
			fmt.Printf("%d. %s\n", i, list)
		}
	},
}

var playlistCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new playlist",
	Long:  `create a new playlist`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		if len(args) > 0 {
			name := strings.Join(args, " ")
			controller.CreatePlayList(name, client)
		} else {
			fmt.Println("no playlist name specified")
		}
	},
}

// remove
var playlistRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a playlist (and its songs)",
	Long:  `remove a playlist (and its songs)`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		controller.RemovePlayList(args[0], client)
	},
}

//add song to a playlist
var playlistAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add music(s) to a playlist",
	Long:  `add music(s) to a playlist`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		if len(args) < 2 {
			fmt.Println("playlist name and query required")
			return
		}
		name := args[0]
		query := strings.Join(args[1:], " ")
		views.SearchThenAddToPlayList(name, query, client)
	},
}

var playlistRemoveSongCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a song from a playlist",
	Long:  `remove a song from a playlist`,
	Run: func(cmd *cobra.Command, args []string) {
		client := controller.GetClient()
		if len(args) < 2 {
			fmt.Println("playlist name and index required")
			return
		}
		name := args[0]
		index, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		controller.RemoveSongFromPlayList(name, index, client)
	},
}
