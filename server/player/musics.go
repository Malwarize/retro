package player

import (
	"path/filepath"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
)

type Music struct {
	Volume *effects.Volume
	Format beep.Format
	Path   string
}

func NewMusic(path string) (Music, error) {
	streamer, format, err := MusicDecode(path)
	if err != nil {
		return Music{}, err
	}
	return Music{
		Volume: &effects.Volume{
			Streamer: streamer,
			Base:     2,
			Silent:   false,
		},
		Format: format,
		Path:   path,
	}, nil
}

func (m Music) String() string {
	return m.Path
}

func (m Music) Name() string {
	return filepath.Base(m.Path)
}

func (m Music) Streamer() beep.StreamSeekCloser {
	return m.Volume.Streamer.(beep.StreamSeekCloser)
}

func (m Music) SetVolume(vp int) {
	if vp == 0 {
		m.Volume.Silent = true
	} else {
		m.Volume.Silent = false
		volume := float64(vp-100) / 16.0
		m.Volume.Volume = volume
	}
}
