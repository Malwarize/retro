package player

import (
	"path/filepath"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
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

func (m Music) DurationN() int {
	speaker.Lock()
	defer speaker.Unlock()
	return m.Streamer().Len()
}

func (m Music) DurationD() time.Duration {
	return m.Format.SampleRate.D(m.DurationN())
}

func (m Music) PositionN() int {
	speaker.Lock()
	defer speaker.Unlock()
	return m.Streamer().Position()
}

func (m Music) PositionD() time.Duration {
	return m.Format.SampleRate.D(m.PositionN())
}

func (m Music) SetPositionN(p int) error { // this indicates where the music play is (samples)
	speaker.Lock()
	defer speaker.Unlock()
	return m.Streamer().Seek(p)
}

func (m Music) SetPositionD(d time.Duration) error {
	dur := m.DurationN()
	new := m.Format.SampleRate.N(d)
	if new < 0 {
		new = 0
	}
	if new > dur {
		new = dur
	}
	return m.SetPositionN(new)
}

func (m Music) Seek(d time.Duration) error {
	return m.SetPositionD(d + m.PositionD())
}
