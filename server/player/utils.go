package player

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
)

type customReadCloser struct {
	io.Reader
	io.Seeker
}

func (crc *customReadCloser) Close() error {
	// Implement the Close method to satisfy the io.Closer interface.
	// This is a no-op for a bytes.Reader.
	return nil
}

// MusicDecode decodes MP3 data from a byte slice and returns a StreamSeekCloser and Format.
func MusicDecode(data []byte) (beep.StreamSeekCloser, beep.Format, error) {
	// Create a bytes.Reader from the data slice for seeking capabilities.
	reader := bytes.NewReader(data)

	// Wrap the reader in the customReadCloser to preserve its seeking capabilities while providing a Close method.
	readerCloser := &customReadCloser{Reader: reader, Seeker: reader}

	// Decode the MP3 data using the custom ReadCloser.
	return mp3.Decode(readerCloser)
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func hash(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func createTmpFile(data []byte) (*os.File, error) {
	f, err := os.CreateTemp("", "goplay_")
	if err != nil {
		return nil, err
	}
	if data != nil {
		f.Write(
			data,
		)
	}
	return f, nil
}
