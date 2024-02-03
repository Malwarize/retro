package views

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"strconv"
	"strings"
	"time"

	"github.com/Malwarize/goplay/client/controller"
	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/bubbles/progress"
)

func reformatDuration(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}

func parseName(name string) string {
	if strings.Contains(name, "/") {
		name = strings.Split(name, "/")[len(strings.Split(name, "/"))-1]
	}
	if strings.Contains(name, shared.Separator) {
		name = strings.Split(name, shared.Separator)[0]
	}
	return name
}

func DisplayStatus(client *rpc.Client) {
	status := controller.GetPlayerStatus(client)
	queue := status.MusicList
	if status.PlayerState == shared.Stopped {
		fmt.Println(stoppedStyle.Render(emojiesStatus[shared.Stopped], " Stopped"))
	} else {

		currentMusicName := queue[status.CurrentMusicIndex]

		currentMusicName = parseName(currentMusicName)
		currentPosition := status.CurrentMusicPosition
		currentPositionStr := reformatDuration(currentPosition)

		totalDurationStr := reformatDuration(status.CurrentMusicLength)

		prog := progress.New(progress.WithSolidFill("#FF0066"))
		prog.SetPercent(0.5)
		prog.ShowPercentage = false
		fmt.Println(progressStyle.Render(prog.ViewAs(currentPosition.Seconds() / status.CurrentMusicLength.Seconds())))

		fmt.Println("   "+playingEmojies[rand.Intn(len(playingEmojies))], currentMusicName)
		fmt.Println(positionStyle.Render(currentPositionStr + " / " + totalDurationStr))

		switch status.PlayerState {
		case shared.Playing:
			fmt.Println(runningStyle.Render(emojiesStatus[shared.Playing], " Playing"))
		case shared.Paused:
			fmt.Println(pausedStyle.Render(emojiesStatus[shared.Paused], " Paused"))
		}
		// display queue
		for i, music := range queue {
			if i == status.CurrentMusicIndex {
				fmt.Println(selectMusicStyle.Render("->", strconv.Itoa(i), ":", parseName(music)))
			} else {
				fmt.Println("  ", i, ":", parseName(music))
			}
		}
	}
	// display tasks
	for target, task := range status.Tasks {
		if task.Error != "" {
			fmt.Println(failedtaskStyle.Render(tasksEmojies[task.Type], " ", target, " ", fmt.Sprintf("%v", task.Error)))
			continue
		}
		switch task.Type {
		case shared.Download:
			fmt.Println(taskStyle.Render(tasksEmojies[task.Type], "Downloading ", target))
		case shared.Search:
			fmt.Println(taskStyle.Render(tasksEmojies[task.Type], "Searching ", target))
		}
	}
}
