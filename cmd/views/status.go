package views

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"strings"
	"time"

	"github.com/Malwarize/goplay/controller"
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
	currentMusicName, _ := shared.ParseCachedFileName(queue[status.CurrentMusicIndex])

	//split by / if exists

	pathSplits := strings.Split(currentMusicName, "/")
	if len(pathSplits) > 1 {
		currentMusicName = pathSplits[len(pathSplits)-1]
	} else {
		currentMusicName = pathSplits[0]
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

}
