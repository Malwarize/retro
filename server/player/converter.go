package player

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Malwarize/retro/config"
)

type Converter struct {
	ffmpegPath  string
	ffprobePath string
}

func NewConverter() (*Converter, error) {
	// check if ffmpegPath and ffprobePath exist and are executable
	ppegPath, err := exec.LookPath(config.GetConfig().PathFFmpeg)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %v", err)
	}

	probePath, err := exec.LookPath(config.GetConfig().PathFFprobe)
	if err != nil {
		return nil, fmt.Errorf("ffprobe not found: %v", err)
	}

	return &Converter{
		ffmpegPath:  ppegPath,
		ffprobePath: probePath,
	}, nil
}

func (c *Converter) ConvertToMP3(inputData []byte) ([]byte, error) {
	cmd := exec.Command(
		c.ffmpegPath,
		"-i", "pipe:0", // Read from stdin
		"-vn",
		"-ar", "44100",
		"-ac", "2",
		"-b:a", "192k",
		"-f", "mp3", // Specify output format
		"pipe:1", // Write to stdout
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = bytes.NewReader(inputData)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error converting to MP3: %v\n%s", err, stderr.String())
	}

	return out.Bytes(), nil
}

func (c *Converter) IsMp3(fileData []byte) (bool, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=format_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-i", "pipe:0", // Read from stdin
	)

	cmd.Stdin = bytes.NewReader(fileData)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf(
			"error running ffprobe to check if fileData is mp3: %v %s",
			err,
			string(out),
		)
	}

	fileFormat := strings.TrimSpace(string(out))
	isMP3 := fileFormat == "mp3"
	return isMP3, nil
}
