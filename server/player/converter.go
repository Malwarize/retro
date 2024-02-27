package player

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Malwarize/goplay/config"
)

type Converter struct {
	ffmpegPath  string
	ffprobePath string
}

func NewConverter() (*Converter, error) {
	// check if ffmpegPath and ffprobePath exist and are executable
	ppegPath, err := exec.LookPath(config.GetConfig().Pathffmpeg)
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found: %v", err)
	}

	probePath, err := exec.LookPath(config.GetConfig().Pathffprobe)
	if err != nil {
		return nil, fmt.Errorf("ffprobe not found: %v", err)
	}

	return &Converter{
		ffmpegPath:  ppegPath,
		ffprobePath: probePath,
	}, nil
}

func (c *Converter) ConvertToMP3(inputData []byte) ([]byte, error) {
	tmpInput, err := createTmpFile(inputData)
	if err != nil {
		return nil, fmt.Errorf("error creating temporary input file: %v", err)
	}
	defer func() {
		tmpInput.Close()
		os.Remove(tmpInput.Name())
	}()

	tmpOutput, err := createTmpFile(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating temporary output file: %v", err)
	}
	defer func() {
		tmpOutput.Close()
		os.Remove(tmpOutput.Name())
	}()

	cmd := exec.Command(
		c.ffmpegPath,
		"-i",
		tmpInput.Name(),
		"-vn",
		"-ar",
		"44100",
		"-ac",
		"2",
		"-b:a",
		"192k",
		tmpOutput.Name(),
		"-y",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error converting to MP3: %v\n%s", err, out)
	}

	data, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading output file: %v", err)
	}

	return data, nil
}

func (c *Converter) IsMp3(fileData []byte) (bool, error) {
	cmd := exec.Command("ffprobe",
		"-v",
		"error",
		"-show_entries",
		"format=format_name",
		"-of",
		"default=noprint_wrappers=1:nokey=1",
	)

	tmpFile, err := createTmpFile(fileData)
	defer os.Remove(tmpFile.Name())
	cmd.Args = append(cmd.Args, tmpFile.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf(
			"error running ffprobe to check if fileData is mp3: %v %v",
			err,
			string(out),
		)
	}
	fileFormat := strings.TrimSpace(string(out))
	isMP3 := fileFormat == "mp3"
	return isMP3, nil
}
