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

func PlayDir(dir string, client *rpc.Client) {
	dir_p, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	entries, err := dir_p.Readdir(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			fmt.Println("Adding music: ", dir+"/"+entry.Name())
			AddMusic(dir+"/"+entry.Name(), client)
		}
	}
	Play(client)
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

func GetPlayerStatus(client *rpc.Client) {
	var reply shared.Status
	err := client.Call("Player.RPCGetPlayerStatus", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Player status:", reply)
}

var client *rpc.Client

func GetClient() *rpc.Client {
	if client == nil {
		var err error
		client, err = rpc.Dial("tcp", "localhost:1234")
		if err != nil {
			fmt.Println("the player " + "localhost:1234" + " not running")
			os.Exit(1)
		}
	}
	return client
}
