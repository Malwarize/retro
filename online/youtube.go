package online

import (
	"context"
	"errors"
	"io"

	"github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	gyoutube "google.golang.org/api/youtube/v3"
)

var apiKey = ""

type youtubeEngine struct {
	Service *gyoutube.Service
	Client  *youtube.Client
}

func newYoutubeEngine() (*youtubeEngine, error) {
	if len(apiKey) == 0 {
		return nil, errors.New("API key not found")
	}
	client := option.WithAPIKey(
		apiKey,
	)
	service, err := gyoutube.NewService(context.Background(), client)
	if err != nil {
		return nil, err
	}

	return &youtubeEngine{
		Service: service,
		Client:  &youtube.Client{},
	}, nil
}

func (yt *youtubeEngine) Search(query string, maxResults int) ([]string, error) {
	call := yt.Service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(int64(maxResults))
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	var videos []string
	for _, item := range response.Items {
		if item.Id.VideoId == "" {
			continue
		}
		videos = append(videos, item.Id.VideoId)
	}
	return videos, nil
}

// returns stream, title, error
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
