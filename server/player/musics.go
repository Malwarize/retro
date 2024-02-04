package player

import (
	"path/filepath"

	"github.com/gopxl/beep"
)

type Music struct {
	Streamer beep.StreamSeekCloser
	Format   beep.Format
	Path     string
}

func NewMusic(path string) (Music, error) {
	streamer, format, err := MusicDecode(path)
	if err != nil {
		return Music{}, err
	}
	return Music{
		Streamer: streamer,
		Format:   format,
		Path:     path,
	}, nil
}

func (m Music) String() string {
	return m.Path
}

func (m Music) Name() string {
	return filepath.Base(m.Path)
}
