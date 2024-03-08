package controller

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/Malwarize/retro/shared"
)

func GetPlayListsNames(client *rpc.Client) []string {
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

func DetectAndAddToPlayList(
	name string,
	query string,
	client *rpc.Client,
) ([]shared.SearchResult, error) {
	args := shared.AddToPlayListArgs{PlayListName: name, Query: query}
	var reply []shared.SearchResult
	err := client.Call("Player.RPCDetectAndAddToPlayList", args, &reply)
	return reply, err
}

func PlayListMusics(name string, client *rpc.Client) []string {
	args := name
	var reply []string
	err := client.Call("Player.RPCPlayListMusics", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func RemoveMusicFromPlayList(name string, indexOrName shared.IntOrString, client *rpc.Client) {
	args := shared.RemoveMusicFromPlayListArgs{
		PlayListName: name,
		IndexOrName:  indexOrName,
	}
	var reply int
	err := client.Call("Player.RPCRemoveMusicFromPlayList", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func PlayListPlayMusic(lname string, indexOrName shared.IntOrString, client *rpc.Client) {
	args := shared.PlayListPlayMusicArgs{
		PlayListName: lname,
		IndexOrName:  indexOrName,
	}
	var reply int
	err := client.Call("Player.RPCPlayListPlayMusic", args, &reply)
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
