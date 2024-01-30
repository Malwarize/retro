package online

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Malwarize/goplay/shared"
)

type CachedFile struct {
	Name  string
	Url   string
	Ftype string // youtube, soundcloud, spotify, ...
}

type CachedFiles struct {
	BaseDir string
	Files   []CachedFile
}

func NewCachedFiles(baseDir string) *CachedFiles {
	return &CachedFiles{
		BaseDir: baseDir,
	}
}

func parseCachedFileName(filename string) (string, string) {
	return strings.Split(filename, "_")[0], strings.Split(filename, "_")[1]
}

func getYoutubeIdFromUrl(url string) string {
	return strings.Split(url, "=")[1]
}

func (cf *CachedFiles) Fetch() error {
	log.Println("Fetching cached files")
	dir := filepath.Join(cf.BaseDir)
	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
		log.Println("Cache dir not found, creating it")
	}

	f, err = os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	ftypes, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, ftype := range ftypes {
		dirPath := filepath.Join(cf.BaseDir, ftype)
		fs, err := os.Open(dirPath)
		if os.IsNotExist(err) {
			err = os.Mkdir(dirPath, 0755)
			if err != nil {
				return err
			}
			log.Println(dirPath, " not found, creating it")
		}

		fs, err = os.Open(dirPath)

		if err != nil {
			return err
		}
		defer fs.Close()
		files, err := fs.Readdirnames(-1)
		if err != nil {
			return err
		}
		for _, file := range files {
			name, url := parseCachedFileName(file)
			cf.Files = append(cf.Files, CachedFile{
				Name:  name,
				Url:   url,
				Ftype: ftype,
			})
		}
	}
	log.Println("Cached files fetched")
	return nil
}

func (cf *CachedFiles) GetFileByUrl(url string, ftype string) (string, error) {
	for _, file := range cf.Files {
		log.Println("Searching for file in cache: ", file.Url, url)
		if file.Url == getYoutubeIdFromUrl(url) && file.Ftype == ftype {
			log.Println("File found in cache: ", url)
			return filepath.Join(cf.BaseDir, ftype, file.Name+"_"+file.Url), nil
		}
	}
	return "", errors.New("file not found")
}

func (cf *CachedFiles) AddFile(filedata []byte, name string, ftype string, url string) {
	cf.Fetch()
	log.Println("Adding file to cache: ", name)
	dirPath := filepath.Join(cf.BaseDir, ftype)
	_, err := os.Open(dirPath)
	if err != nil {
		log.Printf("Error opening dir: %v", err)
		return
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			log.Printf("Error creating dir: %v", err)
			return
		}
		log.Println(dirPath, " not found, creating it")
	}

	name = shared.EscapeSpecialDirChars(name)
	filePath := filepath.Join(dirPath, name+"_"+getYoutubeIdFromUrl(url))
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer f.Close()
	f.Write(filedata)
	cf.Fetch()
}

func (cf *CachedFiles) RemoveFile(name string) error {
	for _, file := range cf.Files {
		if file.Name == name {
			log.Println("Removing file from cache: ", name)
			return os.Remove(filepath.Join(cf.BaseDir, file.Ftype, file.Name+"_"+file.Url))
		}
	}
	return nil
}
