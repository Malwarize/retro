package engines

import (
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Malwarize/retro/logger"
	"github.com/Malwarize/retro/shared"
)

func (yt *youtubeEngine) Exists(videoUrl string) (bool, error) {
	cmd := exec.Command(
		yt.ytdlpPath,
		"--ies",
		"all,-generic",
		videoUrl,
		"--skip-download",
	)
	cmd.Stderr = logger.ERRORLogger.Writer()
	logger.LogInfo("excuting command", cmd.Args)
	// yt-dlp --ies all,-generic https://www.youtube.com/watch?v=videoId
	out, err := cmd.Output()
	if err != nil {
		return false, logger.LogError(
			logger.GError(
				"Check video existence failed",
				err,
			),
			string(out),
		)
	}
	return true, nil
}

func (yt *youtubeEngine) Name() string {
	return "youtube"
}

type youtubeEngine struct {
	ytdlpPath string
}

func NewYoutubeEngine() (*youtubeEngine, error) {
	path := "yt-dlp"
	absPath, err := exec.LookPath(path)
	if err != nil {
		return nil, err
	}
	return &youtubeEngine{
		ytdlpPath: absPath,
	}, nil
}

// Search why I used ytdlp instead of YouTube lib : because ytdlp doesn't need API key to search
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

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(
			line,
			":",
		) {
			dur, err := shared.StringToDuration(line)
			if err != nil {
				continue
			}
			currentResult.Duration = dur
			currentResult.Destination = "https://www.youtube.com/watch?v=" + lines[i-1]
			currentResult.Title = lines[i-2]
			currentResult.Type = yt.Name()
			results = append(results, currentResult)
		}
	}

	return results, nil
}

func (yt *youtubeEngine) getYoutubeTitleFromUrl(url string) (string, error) {
	cmd := exec.Command(
		yt.ytdlpPath, "--get-title", url,
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func giveMeTempFileName() (string, error) {
	tmpSongFile, err := os.CreateTemp("", "retro-youtube-*.mp3")
	defer os.Remove(tmpSongFile.Name())
	if err != nil {
		return "", err
	}
	return tmpSongFile.Name(), nil
}

func (yt *youtubeEngine) Download(videoUrl string) (io.ReadCloser, string, error) {
	videoTitle, err := yt.getYoutubeTitleFromUrl(videoUrl)
	if err != nil {
		return nil, "", err
	}
	logger.LogInfo("Downloading", videoTitle, "from", videoUrl)

	tmpSongFile, err := giveMeTempFileName()
	if err != nil {
		return nil, "", err
	}

	//yt-dlp --extract-audio --audio-format mp3 --output "/tmp/f.mp3" --progress https://www.youtube.com/watch\?v\=-RijT8GW4yw0
	cmd := exec.Command(
		yt.ytdlpPath,
		"--extract-audio",
		"--audio-format",
		"mp3",
		"--no-warning",
		"--output",
		tmpSongFile,
		videoUrl,
	)
	logger.LogInfo("excuting command", cmd.Args)
	cmd.Stderr = logger.ERRORLogger.Writer()
	out, err := cmd.Output()
	if err != nil {
		return nil, "", err
	}

	logger.LogInfo("yt-dlp output:\n", string(out))

	// fill the content in buffer and return it
	buffer, err := os.ReadFile(tmpSongFile)
	if err != nil {
		return nil, "", err
	}

	reader := io.NopCloser(strings.NewReader(string(buffer)))
	// defer func() {
	// 	tmpSongFile.Close()
	// 	logger.LogInfo(
	// 		"Removing tmp file", tmpSongFile.Name(),
	// 	)
	// 	os.Remove(tmpSongFile.Name())
	// }()

	return reader, videoTitle, nil
}

func (yt *youtubeEngine) MaxResults() int {
	return 10
}
