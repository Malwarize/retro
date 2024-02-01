package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Malwarize/goplay/cmd/views"
	"github.com/Malwarize/goplay/controller"
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

		if file, err := cmd.Flags().GetString("file"); err == nil && file != "" {
			controller.PlayFile(file, client)
			os.Exit(0)
		}

		if youtube, err := cmd.Flags().GetString("youtube"); err == nil && youtube != "" {
			controller.PlayYoutube(youtube, client)
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
		controller.GetPlayerStatus(client)
	},
}
