package player

import (
	"os"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

func MusicDecode(path string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, err
	}
	return mp3.Decode(f)
}
