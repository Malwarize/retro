package views

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Malwarize/goplay/client/controller"
	"github.com/Malwarize/goplay/config"
	"github.com/Malwarize/goplay/shared"
	"github.com/charmbracelet/bubbles/progress"
)

func reformatDuration(duration time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), int(duration.Seconds())%60)
}

func parseName(name string) string {
	name = filepath.Base(name)
	if strings.Contains(name, config.GetConfig().Separator) {
		name = strings.Split(name, config.GetConfig().Separator)[0]
	}
	return name
}
func convertVolumeToEmojie(volume int) string {
	if volume == 0 {
		return volumeLevels[0]
	}
	if volume < 50 {
		return volumeLevels[1]
	}
	if volume < 85 {
		return volumeLevels[2]
	}
	return volumeLevels[3]
}

func DisplayStatus(client *rpc.Client) {
	status := controller.GetPlayerStatus(client)
	queue := status.MusicQueue
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
			fmt.Println(runningStyle.Render(emojiesStatus[shared.Playing], " Playing", convertVolumeToEmojie(status.Volume)))
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
			switch task.Type {
			case shared.Downloading:
				fmt.Println(failedtaskStyle.Render(failedEmojie, "Failed to download ", target, ":", task.Error))
			case shared.Searching:
				fmt.Println(failedtaskStyle.Render(failedEmojie, "Failed to search ", target, ":", task.Error))
			default:
				fmt.Println(failedtaskStyle.Render(failedEmojie, "Failed to ", target, ":", task.Error))
			}
			continue
		}
		switch task.Type {
		case shared.Downloading:
			fmt.Println(taskStyle.Render(tasksEmojies[task.Type], "Downloading ", target))
		case shared.Searching:
			fmt.Println(taskStyle.Render(tasksEmojies[task.Type], "Searching ", target))
		}
	}
}
