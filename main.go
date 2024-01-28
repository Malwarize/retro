package main

import (
	"fmt"
	"os"

	"github.com/Malwarize/goplay/cmd"
	"github.com/Malwarize/goplay/player"
	"github.com/gopxl/beep/mp3"
)

func getMusic(musicPath string) player.Music {
	f, err := os.Open(musicPath)
	if err != nil {
		panic(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		panic(err)
	}
	return player.NewMusic(musicPath, streamer, format)
}

func tester() {
	p := player.GetPlayer()
	music_path := "audio/00.mp3"
	music := getMusic(music_path)
	p.AddMusic(music)
	music_path = "audio/01.mp3"
	music = getMusic(music_path)
	p.AddMusic(music)
	fmt.Println("playing music: ", p.GetCurrentMusic())
	fmt.Println("Duration: ", p.GetCurrentMusicLength())

	p.Play()
	p.Wait()
}

func main() {
	cmd.Execute()
}
