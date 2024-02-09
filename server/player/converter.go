package player

import (
	"fmt"
	"log"
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

func (c *Converter) ConvertToMP3(inputFile string) error {
	tmpFile, err := os.CreateTemp("", "goplay")
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	outputFile := tmpFile.Name() + ".mp3"
	defer os.Remove(outputFile)

	cmd := exec.Command(c.ffmpegPath, "-i", inputFile, "-vn", "-ar", "44100", "-ac", "2", "-b:a", "192k", outputFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("error converting to MP3:", string(out))
		return fmt.Errorf("error converting to MP3: %v", err)
	}

	err = copyFile(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	err = os.Remove(outputFile)
	if err != nil {
		return fmt.Errorf("error converting to MP3: %v", err)
	}
	return nil
}

func (c *Converter) IsMp3(file string) (bool, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=format_name", "-of", "default=noprint_wrappers=1:nokey=1", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error checking file format:", string(output))
		return false, fmt.Errorf("error checking file format: %v", err)
	}
	fileFormat := strings.TrimSpace(string(output))
	isMP3 := strings.EqualFold(fileFormat, "mp3")
	return isMP3, nil
}
