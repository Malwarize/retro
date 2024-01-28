package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goplay",
	Short: "goplay is a music player",
	Long:  `goplay is a music player`,
	Run: func(_ *cobra.Command, _ []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	playCmd.Flags().StringP("dir", "d", "", "add all songs in a directory to the queue")
	playCmd.Flags().StringP("file", "f", "", "play a single song")

	serverCmd.Flags().StringP("port", "p", "1234", "port to listen on")

	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(seekCmd)
	rootCmd.AddCommand(statusCmd)

}
