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

func DisplayStatus(client *rpc.Client) {
	status := controller.GetPlayerStatus(client)

	if status.PlayerState == shared.Stopped {
		fmt.Println(stoppedStyle.Render(emojiesStatus[shared.Stopped], " Stopped"))
		return
	}
	queue := status.MusicList
	currentMusicName := queue[status.CurrentMusicIndex]

	//split by / if exists
	if strings.Contains(currentMusicName, "/") {
		currentMusicName = strings.Split(currentMusicName, "/")[len(strings.Split(currentMusicName, "/"))-1]
		if strings.Contains(currentMusicName, shared.Separator) {
			currentMusicName = strings.Split(currentMusicName, shared.Separator)[0]
		}
	}

	currentPosition := status.CurrentMusicPosition
	currentPositionStr := reformatDuration(currentPosition)

	totalDurationStr := reformatDuration(status.CurrentMusicLength)

	prog := progress.New(progress.WithSolidFill("#FF0066"))
	prog.SetPercent(0.5)
	prog.ShowPercentage = false
	fmt.Println(progressStyle.Render(prog.ViewAs(currentPosition.Seconds() / status.CurrentMusicLength.Seconds())))

	fmt.Println(playingEmojies[rand.Intn(len(playingEmojies))], " ", currentMusicName)
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
			fmt.Println(selectMusicStyle.Render("->", strconv.Itoa(i), ":", music))
		} else {
			fmt.Println("  ", i, ":", music)
		}
	}
}
