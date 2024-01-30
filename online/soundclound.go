package online

// import (
// 	"compress/gzip"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/http/cookiejar"
// 	"strings"
// )

// var endpoint = "https://soundcloud.com/search/sounds?q="

// type SoundCloudClient struct {
// 	client *http.Client
// }

// var (
// 	SoundCloudClientInstance *SoundCloudClient
// )

// func GetSoundCloudClient() *SoundCloudClient {
// 	if SoundCloudClientInstance == nil {
// 		SoundCloudClientInstance = NewSoundCloudClient()
// 	}

// 	return SoundCloudClientInstance
// }

// func NewSoundCloudClient() *SoundCloudClient {
// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	httpsClient := &http.Client{
// 		Jar: jar,
// 	}

// 	return &SoundCloudClient{
// 		client: httpsClient,
// 	}
// }

// var GET_HEADERS = map[string]string{
// 	"Accept":          "*/*",
// 	"Accept-Encoding": "gzip, deflate, br",
// 	"Accept-Language": "en-US,en;q=0.9",
// 	"Connection":      "keep-alive",
// 	"Host":            "soundcloud.com",
// 	"Referer":         "https://soundcloud.com/",
// 	"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
// }

// type SoundCloudMusic struct {
// 	Title      string
// 	Link       string
// 	StreamLink string
// }

// func (sc *SoundCloudClient) HttpDo(req *http.Request) (*http.Response, error) {
// 	for key, value := range GET_HEADERS {
// 		req.Header.Set(key, value)
// 	}
// 	return sc.client.Do(req)
// }

// func (sc *SoundCloudClient) SearchForMusic(query string) ([]SoundCloudMusic, error) {
// 	url := endpoint + query

// 	var musics []SoundCloudMusic
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := sc.HttpDo(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var reader io.Reader
// 	switch resp.Header.Get("Content-Encoding") {
// 	case "gzip":
// 		reader, err = gzip.NewReader(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer reader.(*gzip.Reader).Close()
// 	default:
// 		reader = resp.Body
// 	}

// 	if err != nil {
// 		return nil, err
// 	}
// 	body, err := io.ReadAll(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("status code %d", resp.StatusCode)
// 	}
// 	matches, err := parseSoundCloudMusicUrlFromHtml(string(body))
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, match := range matches {
// 		musics = append(musics, SoundCloudMusic{
// 			Title: strings.Split(match, "/")[len(strings.Split(match, "/"))-1],
// 			Link:  match,
// 		})
// 	}
// 	return musics, nil
// }

// func (sc *SoundCloudClient) FetchMusicData(music *SoundCloudMusic) error {
// 	req, err := http.NewRequest("GET", music.Link, nil)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := sc.HttpDo(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}
// 	var reader io.Reader
// 	switch resp.Header.Get("Content-Encoding") {
// 	case "gzip":
// 		reader, err = gzip.NewReader(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		defer reader.(*gzip.Reader).Close()
// 	default:
// 		reader = resp.Body
// 	}

// 	if err != nil {
// 		return err
// 	}
// 	body, err := io.ReadAll(reader)
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode != 200 {
// 		return fmt.Errorf("status code %d", resp.StatusCode)
// 	}
// 	title, err := parseSoundCloudMusicTitleFromHtml(string(body))
// 	if err != nil {
// 		return err
// 	}
// 	music.Title = title

// 	//
// 	return nil
// }
