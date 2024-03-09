package controller

import (
	"fmt"
	"net/rpc"
	"os"
	"strings"
)

func GetLogs(client *rpc.Client) []string {
	var reply []string
	err := client.Call("Player.GetLogs", 0, &reply)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return reply
}

// prettify the log output add emoji and color

const (
	infoEmojie  = "📄"
	warnEmojie  = "⚠️"
	errorEmojie = "❌"
)

func prettifyLogs(logs []string) []string {
	var prettyLogs []string
	for _, log := range logs {
		if strings.Contains(log, "ERROR:") {
			log = strings.ReplaceAll(log, "ERROR:", errorEmojie)
		}
		if strings.Contains(log, "WARN:") {
			log = strings.ReplaceAll(log, "WARN:", warnEmojie)
		}
		if strings.Contains(log, "INFO:") {
			log = strings.ReplaceAll(log, "INFO:", infoEmojie)
		}
		prettyLogs = append(prettyLogs, log)
	}
	return prettyLogs
}

func PrintErrorLogs(client *rpc.Client) {
	logs := GetLogs(client)
	logs = prettifyLogs(logs)
	for _, log := range logs {
		if strings.Contains(log, errorEmojie) {
			fmt.Println(log)
		}
	}
}

func PrintInfoLogs(client *rpc.Client) {
	logs := GetLogs(client)
	logs = prettifyLogs(logs)
	for _, log := range logs {
		if strings.Contains(log, infoEmojie) {
			fmt.Println(log)
		}
	}
}

func PrintWarnLogs(client *rpc.Client) {
	logs := GetLogs(client)
	logs = prettifyLogs(logs)
	for _, log := range logs {
		if strings.Contains(log, warnEmojie) {
			fmt.Println(log)
		}
	}
}

func PrintAllLogs(client *rpc.Client) {
	logs := GetLogs(client)
	logs = prettifyLogs(logs)
	for _, log := range logs {
		fmt.Println(log)
	}
}
