package cmd

import (
	"fmt"
	"github.com/Malwarize/retro/client/controller"
	"github.com/Malwarize/retro/shared"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "retro",
	Short: "retro is a music player",
	Long: `retro is a music player
retro is client for retroPlayer server
you can control retroPlayer server like any other systemd service
retro [command] --help for more information about a command`,
	Version: shared.Version,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("retro is a music player")
		fmt.Println("use retro --help to see available commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	client, err := controller.GetClient()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	if controller.IsUpdatePromptEnabled(client) && controller.IsUpdateAvailable(client) {
		fmt.Print("A new version is available, do you want to update? [m]mute this prompt [Y/n]")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return
		}
		if input == "y" || input == "Y" {
			controller.Update(client)
			os.Exit(0)
		} else if input == "m" {
			controller.DisableTheUpdatePrompt(client)
		}
	}

	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(seekCmd)
	rootCmd.AddCommand(seekBackCmd)
	rootCmd.AddCommand(volumeCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(setThemeCmd)
	rootCmd.AddCommand(playlistCmd)
	rootCmd.AddCommand(logCmd)

	playlistCmd.AddCommand(playlistCreateCmd)
	playlistCmd.AddCommand(playlistRemoveCmd)
	playlistCmd.AddCommand(playlistAddCmd)
	playlistCmd.AddCommand(playlistPlayCmd)

	logCmd.AddCommand(logErrCmd)
	logCmd.AddCommand(logInfoCmd)
	logCmd.AddCommand(logWarnCmd)

	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cleanCacheCmd)
}
