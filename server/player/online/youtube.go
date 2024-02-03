package online

import (
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"

	"github.com/Malwarize/goplay/shared"
)

// import (
// 	"context"
// 	"errors"
// 	"io"
// 	"strings"
// 	"time"

// 	"github.com/Malwarize/goplay/shared"
// 	"github.com/kkdai/youtube/v2"
// 	"google.golang.org/api/option"
// 	gyoutube "google.golang.org/api/youtube/v3"
// )

// var apiKey = ""

// type youtubeEngine struct {
// 	Service *gyoutube.Service
// 	Client  *youtube.Client
// }

// func getYoutubeIdFromUrl(url string) string {
// 	return strings.Split(url, "=")[1]
// }

// func newYoutubeEngine() (*youtubeEngine, error) {
// 	if len(apiKey) == 0 {
// 		return nil, errors.New("API key not found")
// 	}
// 	client := option.WithAPIKey(
// 		apiKey,
// 	)
// 	service, err := gyoutube.NewService(context.Background(), client)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &youtubeEngine{
// 		Service: service,
// 		Client:  &youtube.Client{},
// 	}, nil
// }

// var DEBUG = false

// func (yt *youtubeEngine) Search(query string, maxResults int) ([]shared.SearchResult, error) {
// 	if DEBUG {
// 		time.Sleep(1 * time.Second)
// 		return []shared.SearchResult{
// 			{
// 				Title:       "Test",
// 				Destination: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
// 				Type:        "youtube",
// 			},
// 			{
// 				Title:       "Test2",
// 				Destination: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
// 				Type:        "youtube",
// 			},
// 			{
// 				Title:       "Test3",
// 				Destination: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
// 				Type:        "youtube",
// 			},
// 		}, nil
// 	}
// 	call := yt.Service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(int64(maxResults))
// 	response, err := call.Do()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var videos []shared.SearchResult
// 	for _, item := range response.Items {
// 		if item.Id.VideoId == "" {
// 			continue
// 		}
// 		videos = append(videos, shared.SearchResult{
// 			Title:       item.Snippet.Title,
// 			Destination: "https://www.youtube.com/watch?v=" + item.Id.VideoId,
// 			Type:        "youtube",
// 		})
// 	}
// 	return videos, nil
// }

// // returns stream, title, error
func (yt *youtubeEngine) Download(videoUrl string) (io.ReadCloser, string, error) {

}

func (yt *youtubeEngine) Exists(videoUrl string) (bool, error) {
	_, err := yt.Client.GetVideo(videoUrl)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (yt *youtubeEngine) GetName() string {
	return "youtube"
}

type youtubeEngine struct {
	ytdlpPath string
	Client    *youtube.Client
}

func newYoutubeEngine() (*youtubeEngine, error) {
	// check if yt-dlp exits\
	path := "yt-dlp"
	absPath, err := exec.LookPath(path)
	if err != nil {
		return nil, err
	}
	client := &youtube.Client{}
	return &youtubeEngine{
		ytdlpPath: absPath,
		Client:    client,
	}, nil
}

func (yt *youtubeEngine) Search(query string, maxResults int) ([]shared.SearchResult, error) {
	cmd := exec.Command(yt.ytdlpPath, "--get-id", "--get-title", "--skip-download", "--flat-playlist", "ytsearch"+strconv.Itoa(maxResults)+":"+query)
	// --get-id --get-title --skip-download --simulate ytsearch<maxResults>:<query>
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var results []shared.SearchResult
	for i := 0; i < len(lines)-1; i += 2 {
		results = append(results, shared.SearchResult{
			Title:       lines[i],
			Destination: "https://www.youtube.com/watch?v=" + lines[i+1],
			Type:        "youtube",
		})
	}
	return results, nil
}
