package controller

import (
	"fmt"
	"net/rpc"
	"os"
)

//func (p *Player) IsUpdateAvailable(_ int, reply *bool) error {
//	logger.LogInfo("NeedUpdate called")
//	*reply = p.up.IsUpdateAvailable
//	logger.LogInfo("NeedUpdate done with reply :", reply)
//	return nil
//}
//
//func (p *Player) Update(_ int, reply *int) error {
//	logger.LogInfo("Update called")
//	err := p.up.Update()
//	logger.LogInfo("Update done")
//	return err
//}
//
//func (p *Player) DisableTheUpdatePrompt(_ int, reply *int) error {
//	logger.LogInfo("DisableTheUpdatePrompt called")
//	p.up.EnableUpdatePrompt = false
//	*reply = 1
//	logger.LogInfo("DisableTheUpdatePrompt done")
//	return nil
//}
//
//func (p *Player) IsUpdatePromptEnabled(_ int, reply *bool) error {
//	logger.LogInfo("IsUpdatePromptEnabled called")
//	*reply = p.up.EnableUpdatePrompt
//	logger.LogInfo("IsUpdatePromptEnabled done with reply :", reply)
//	return nil
//}

func IsUpdateAvailable(client *rpc.Client) bool {
	var reply bool
	err := client.Call("Player.IsUpdateAvailable", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

func Update(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.Update", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DisableTheUpdatePrompt(client *rpc.Client) {
	args := 0
	var reply int
	err := client.Call("Player.DisableTheUpdatePrompt", args, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func IsUpdatePromptEnabled(client *rpc.Client) bool {
	var reply bool
	err := client.Call("Player.IsUpdatePromptEnabled", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}
