package online

import (
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kkdai/youtube/v2"

	"github.com/Malwarize/goplay/logger"
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
// func (yt *youtubeEngine) Download(videoUrl string) (io.ReadCloser, string, error) {
// 	tmpFile, err := os.CreateTemp("", "goplay-youtube-*.mp3")
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer tmpFile.Close()

// 	// --audio-format mp3 --get-title -o tmp.mp3
// 	cmd := exec.Command(yt.ytdlpPath, "-x", "--audio-format", "mp3", "--get-title", "-o", tmpFile.Name(), videoUrl)
// 	out, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return nil, "", err
// 	}

// 	title := strings.TrimSpace(string(out))
// 	tmpReader, err := os.Open(tmpFile.Name())
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	var buffer []byte
// 	buffer, err = io.ReadAll(tmpReader)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	reader := io.NopCloser(strings.NewReader(string(buffer)))
// 	os.Remove(tmpFile.Name())
// 	return reader, title, nil
// }

func (yt *youtubeEngine) Exists(videoUrl string) (bool, error) {
	_, err := yt.Client.GetVideo(videoUrl)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (yt *youtubeEngine) Name() string {
	return "youtube"
}

type youtubeEngine struct {
	ytdlpPath string
	Client    *youtube.Client
}

func newYoutubeEngine() (*youtubeEngine, error) {
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

// why I used ytdlp instead of youtube lib : because ytdlp doesn't need API key to search
func (yt *youtubeEngine) Search(query string, maxResults int) ([]shared.SearchResult, error) {
	cmd := exec.Command(
		yt.ytdlpPath,
		"--get-id",
		"--get-title",
		"--get-duration",
		"--skip-download",
		"--flat-playlist",
		"ytsearch"+strconv.Itoa(maxResults)+":"+query,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	logger.LogInfo("yt-dlp output:\n", string(out))
	lines := strings.Split(string(out), "\n")
	// remove last empty line
	if len(lines) < 1 {
		logger.LogInfo("Invalid yt-dlp output, len(lines) is 0")
	}
	lines = lines[:len(lines)-1]

	var results []shared.SearchResult
	var currentResult shared.SearchResult
	if len(lines)%3 != 0 {
		logger.LogInfo("Invalid yt-dlp output, len(lines) is not devideable by 3 its ", len(lines))
		return nil, nil
	}

	for i := 0; i < len(lines); i += 3 {
		currentResult.Title = lines[i]
		id := lines[i+1]
		if len(id) != 11 {
			continue
		}
		currentResult.Destination = "https://www.youtube.com/watch?v=" + lines[i+1]
		sdur := lines[i+2]
		if dur, err := shared.StringToDuration(sdur); err == nil {
			currentResult.Duration = dur
		} else {
			continue
		}
		results = append(results, currentResult)

	}
	if currentResult.Destination != "" {
		results = append(results, currentResult)
	}

	return results, nil
}

// returns stream, title, error
// why don't use ytdlp : because this lib much faster
func (yt *youtubeEngine) Download(videoUrl string) (io.ReadCloser, string, error) {
	video, err := yt.Client.GetVideo(videoUrl)
	if err != nil {
		return nil, "", err
	}
	formats := video.Formats.Itag(140)
	stream, _, err := yt.Client.GetStream(video, &formats[0])
	if err != nil {
		return nil, "", err
	}
	return stream, video.Title, nil
}
