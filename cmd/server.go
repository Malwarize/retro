package cmd

import (
	"fmt"

	"github.com/Malwarize/goplay/player"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Long:  `start the server`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Println(err)
			return
		}
		player.StartIPCServer(port)
	},
}
