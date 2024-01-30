package online

import (
	"errors"
	"os"
	"strings"

	"log"
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

func (cf *CachedFiles) Fetch() error {
	log.Println("Fetching cached files")
	f, err := os.Open(cf.BaseDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(cf.BaseDir, 0755)
		if err != nil {
			return err
		}
		log.Println("Cache dir not found, creating it")
	}

	f, err = os.Open(cf.BaseDir)
	if err != nil {
		return err
	}
	defer f.Close()
	ftypes, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, ftype := range ftypes {
		fs, err := os.Open(cf.BaseDir + "/" + ftype)
		if os.IsNotExist(err) {
			err = os.Mkdir(cf.BaseDir+"/"+ftype, 0755)
			if err != nil {
				return err
			}
			log.Println(cf.BaseDir + "/" + ftype + " not found, creating it")
		}

		fs, err = os.Open(cf.BaseDir + "/" + ftype)

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
	return nil
}

func (cf *CachedFiles) GetFileByUrl(url string, ftype string) (string, error) {
	for _, file := range cf.Files {
		if file.Url == url && file.Ftype == ftype {
			log.Println("File found in cache: ", url)
			return cf.BaseDir + "/" + ftype + "/" + file.Name + "_" + file.Url, nil
		}
	}
	return "", errors.New("file not found")
}

func (cf *CachedFiles) GetFileByName(name string) (string, error) {
	for _, file := range cf.Files {
		if file.Name == name {
			log.Println("File found in cache: ", name)
			return cf.BaseDir + "/" + file.Ftype + "/" + file.Name + "_" + file.Url, nil
		}
	}
	return "", errors.New("file not found")
}

func (cf *CachedFiles) AddFile(filedata []byte, name string, ftype string, url string) {
	cf.Fetch()
	log.Println("Adding file to cache: ", name)
	f, err := os.Open(cf.BaseDir + "/" + ftype)
	if os.IsNotExist(err) {
		err = os.Mkdir(cf.BaseDir+"/"+ftype, 0755)
		if err != nil {
			log.Printf("Error creating dir: %v", err)
			return
		}
		log.Println(cf.BaseDir + "/" + ftype + " not found, creating it")
	}

	f, err = os.Create(cf.BaseDir + "/" + ftype + "/" + name + "_" + url)
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
			return os.Remove(cf.BaseDir + "/" + file.Ftype + "/" + file.Name + "_" + file.Url)
		}
	}
	return nil
}
