package controller

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/Malwarize/goplay/shared"
)

func PlayListsNames(client *rpc.Client) []string {
	var reply []string
	err := client.Call("Player.RPCPlayListsNames", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func CreatePlayList(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCCreatePlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RemovePlayList(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCRemovePlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DetectAndAddToPlayList(name string, query string, client *rpc.Client) []shared.SearchResult {
	args := shared.AddToPlayListArgs{PlayListName: name, Query: query}
	var reply []shared.SearchResult
	err := client.Call("Player.RPCDetectAndAddToPlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func PlayListSongs(name string, client *rpc.Client) []string {
	args := name
	var reply []string
	err := client.Call("Player.RPCPlayListSongs", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func RemoveSongFromPlayList(name string, index int, client *rpc.Client) {
	args := shared.RemoveSongFromPlayListArgs{PlayListName: name, Index: index}
	var reply int
	err := client.Call("Player.RPCRemoveSongFromPlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayListPlaySong(lname string, index int, client *rpc.Client) {
	args := shared.PlayListPlaySongArgs{PlayListName: lname, Index: index}
	var reply int
	err := client.Call("Player.RPCPlayListPlaySong", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayListPlayAll(name string, client *rpc.Client) {
	args := name
	var reply int
	err := client.Call("Player.RPCPlayListPlayAll", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
