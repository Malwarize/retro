package controller

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/Malwarize/goplay/shared"
)

func Play(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCPlay", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayFile(file string, client *rpc.Client) {
	fmt.Println("Adding music: ", file)
	AddMusic(file, client)
	Play(client)
}

func Next(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCNext", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Prev(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCPrev", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Pause(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCPause", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Resume(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCResume", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Stop(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.RPCStop", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayDir(dir string, client *rpc.Client) {
	args := dir
	var reply int
	err := client.Call("Player.RPCAddDir", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func AddMusic(music string, client *rpc.Client) {
	args := music
	var reply int
	err := client.Call("Player.RPCAddMusic", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Seek(d int, client *rpc.Client) {
	args := d
	var reply int
	err := client.Call("Player.RPCSeek", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Volume(vp int, client *rpc.Client) {
	if vp > 100 {
		// health warning
		fmt.Print("    ⚠️ Volume greater than 100% may damage your ears, skip this warning? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			return
		}
	}
	args := vp
	var reply int
	err := client.Call("Player.RPCVolume", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Remove(index int, client *rpc.Client) {
	args := index
	var reply int
	err := client.Call("Player.RPCRemoveMusic", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetPlayerStatus(client *rpc.Client) shared.Status {
	var reply shared.Status
	err := client.Call("Player.RPCGetPlayerStatus", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func PlayYoutube(url string, client *rpc.Client) {
	args := url
	var reply int
	err := client.Call("Player.RPCPlayYoutube", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("this may take a while before the music starts playing")
}

func DetectAndPlay(query string, client *rpc.Client) []shared.SearchResult {
	var reply []shared.SearchResult
	err := client.Call("Player.RPCDetectAndPlay", query, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func GetTheme(client *rpc.Client) string {
	var reply string
	err := client.Call("Player.RPCGetTheme", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func SetTheme(theme string, client *rpc.Client) {
	args := theme
	var reply int
	err := client.Call("Player.RPCSetTheme", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var client *rpc.Client

func GetClient() *rpc.Client {
	if client == nil {
		var err error
		client, err = rpc.Dial("tcp", "localhost:3131")
		if err != nil {
			fmt.Println("the player " + "localhost:3131" + " not running")
			os.Exit(1)
		}
	}
	return client
}
