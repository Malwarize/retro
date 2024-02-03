package online

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Malwarize/goplay/shared"
)

type CachedFile struct {
	Name  string
	Key   string
	Ftype string // youtube, soundcloud, spotify, ...online
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

func sanitizeName(name string) string {
	re := regexp.MustCompile(`[\/\\\:\*\?\"\<\>\|]`)
	return re.ReplaceAllString(name, "")
}

// check if item already exists in cache

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
		clear(cf.Files)
		for _, file := range files {
			name, key := shared.ParseCachedFileName(file)
			cf.Files = append(cf.Files, CachedFile{
				Name:  name,
				Key:   key,
				Ftype: ftype,
			})
		}
	}
	log.Println("Cached files fetched")
	return nil
}

func (cf *CachedFiles) GetFileByKey(key string, ftype string) (string, error) {
	for _, file := range cf.Files {
		if file.Key == sanitizeName(key) && file.Ftype == ftype {
			log.Println("File found in cache: ", key)
			return filepath.Join(cf.BaseDir, ftype, shared.CombineNameWithKey(file.Name, file.Key)), nil
		}
	}
	return "", errors.New("file not found")
}

func (cf *CachedFiles) GetFileByName(name string, ftype string) (string, error) {
	for _, file := range cf.Files {
		if file.Name == name && file.Ftype == ftype {
			log.Println("File found in cache: ", name)
			return filepath.Join(cf.BaseDir, ftype, shared.CombineNameWithKey(file.Name, file.Key)), nil
		}
	}
	return "", errors.New("file not found")
}

func (cf *CachedFiles) Search(query string) []string {
	var results []string
	for _, file := range cf.Files {
		if strings.Contains(strings.ToLower(file.Name), strings.ToLower(sanitizeName(query))) {
			log.Println("File found in cache: ", file.Name)
			results = append(results, filepath.Join(cf.BaseDir, file.Ftype, shared.CombineNameWithKey(file.Name, file.Key)))
		}
	}
	return results
}

func (cf *CachedFiles) AddFile(filedata []byte, name string, ftype string, key string) string {
	cf.Fetch()
	log.Println("Adding file to cache: ", name)
	dirPath := filepath.Join(cf.BaseDir, ftype)
	_, err := os.Open(dirPath)
	if err != nil {
		log.Printf("Error opening dir: %v", err)
		return ""
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			log.Printf("Error creating dir: %v", err)
			return ""
		}
		log.Println(dirPath, "not found, creating it")
	}

	filePath := filepath.Join(dirPath, sanitizeName(shared.CombineNameWithKey(name, key)))
	log.Println("Writing file to: ", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return ""
	}
	defer f.Close()
	f.Write(filedata)
	cf.Fetch()
	return filePath
}
