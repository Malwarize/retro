package player

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Malwarize/goplay/server/player/online"
	"github.com/Malwarize/goplay/shared"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

var PlayerInstance *Player

type Music struct {
	Name     string
	Streamer beep.StreamSeekCloser
	Format   beep.Format
}

func (m Music) String() string {
	return m.Name
}

const (
	Playing = iota
	Paused
	Stopped
)

type Player struct {
	MusicList         []Music
	CurrentMusicIndex int
	playerState       int
	done              chan struct{}
	initialised       bool
	Converter         *Converter
	Director          *online.OnlineDirector
	Tasks             map[string]shared.Task
	mu                sync.Mutex
}

func NewMusic(name string, streamer beep.StreamSeekCloser, format beep.Format) Music {
	return Music{
		Name:     name,
		Streamer: streamer,
		Format:   format,
	}
}

func NewPlayer() *Player {
	converter, err := NewConverter()
	if err != nil {
		log.Fatal(err)
	}
	director, err := online.NewDefaultDirector()
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		MusicList:         make([]Music, 0),
		CurrentMusicIndex: 0,
		playerState:       Stopped,
		done:              make(chan struct{}),
		initialised:       false,
		Converter:         converter,
		Director:          director,
		Tasks:             make(map[string]shared.Task),
	}
}

func (p *Player) addTask(target string, typeTask int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Tasks[target] = shared.Task{
		Type:  typeTask,
		Error: "",
	}
}

func (p *Player) removeTask(target string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Tasks, target)
}

func (p *Player) errorifyTask(target string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	task, ok := p.Tasks[target]
	if ok {
		task.Error = err.Error()
		p.Tasks[target] = task
	}
}

func GetPlayer() *Player {
	if PlayerInstance == nil {
		PlayerInstance = NewPlayer()
	}
	return PlayerInstance
}

func (p *Player) AddMusic(music Music) {
	p.MusicList = append(p.MusicList, music)
}

func (p *Player) Play() {
	music := p.MusicList[p.CurrentMusicIndex]
	if !p.initialised {
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
		p.initialised = true
	} else {
		speaker.Clear()
		speaker.Init(music.Format.SampleRate, music.Format.SampleRate.N(time.Second/10))
	}
	p.playerState = Playing
	go func() {
		done := make(chan struct{})
		speaker.Play(beep.Seq(music.Streamer, beep.Callback(func() {
			done <- struct{}{}
		})))
		<-done
		p.Next()
	}()
}

func (p *Player) AddMusicFromFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.AddMusic(NewMusic(path, streamer, format))
	if p.playerState == Stopped {
		p.Play()
	}
}

func (p *Player) youtubeToMusic(urlOrId string) (Music, error) {
	reader, path, err := p.Director.Download("youtube", urlOrId)
	if err != nil {
		p.errorifyTask(urlOrId, err)
		return Music{}, err
	}
	isMp3, err := p.Converter.IsMp3(path)
	if err != nil {
		p.errorifyTask(urlOrId, err)
		return Music{}, err
	}
	if len(strings.Split(path, "/")) < 1 {
		return Music{}, errors.New("Music name not found")
	}
	musicName := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	if !isMp3 {
		err = p.Converter.ConvertToMP3(path)
		if err != nil {
			p.errorifyTask(urlOrId, err)
			return Music{}, err
		}
		reader, err = os.Open(path)
		if err != nil {
			p.errorifyTask(urlOrId, err)
			return Music{}, err
		}
		streamer, format, err := mp3.Decode(reader)
		if err != nil {
			p.errorifyTask(urlOrId, err)
			return Music{}, err
		}
		return NewMusic(musicName, streamer, format), nil
	}
	streamer, format, err := mp3.Decode(reader)
	if err != nil {
		p.errorifyTask(urlOrId, err)
		return Music{}, err
	}
	return NewMusic(musicName, streamer, format), nil
}

func (p *Player) AddMusicFromYoutube(urlOrQueryOrID string) {
	music, err := p.youtubeToMusic(urlOrQueryOrID)
	if err != nil {
		fmt.Println(err)
		return
	}
	p.AddMusic(music)
	if p.playerState == Stopped {
		p.Play()
	}
}

func (p *Player) convertAndPlayInTemp(path string) bool {
	f, err := os.CreateTemp("", "goplay")
	defer os.Remove(f.Name())
	if err != nil {
		log.Println(err)
		return false
	}
	sourceFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = io.Copy(f, sourceFile)
	if err != nil {
		log.Println(err)
	}

	err = p.Converter.ConvertToMP3(f.Name())
	if err != nil {
		log.Println(err)
		return false
	}

	p.AddMusicFromFile(f.Name())
	return true
}

func (p *Player) AddMusicsFromDir(dirPath string) {
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	entries, err := dir.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			isMp3, err := p.Converter.IsMp3(dirPath + "/" + entry.Name())
			if err != nil {
				log.Println(err)
			}
			if isMp3 {
				p.AddMusicFromFile(dirPath + "/" + entry.Name())
			} else {
				p.convertAndPlayInTemp(dirPath + "/" + entry.Name())
			}
		}
	}
	if p.playerState == Stopped {
		p.Play()
	}
}
func (p *Player) Next() {
	if len(p.MusicList) < 1 || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.MusicList[p.CurrentMusicIndex].Streamer.Seek(0)
	p.CurrentMusicIndex++
	if p.CurrentMusicIndex >= len(p.MusicList) {
		p.CurrentMusicIndex = 0
	}
	p.Play()
}

func (p *Player) Prev() {
	if len(p.MusicList) < 1 || p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	p.MusicList[p.CurrentMusicIndex].Streamer.Seek(0)
	p.CurrentMusicIndex--
	if p.CurrentMusicIndex < 0 {
		p.CurrentMusicIndex = len(p.MusicList) - 1
	}
	p.Play()
}

func (p *Player) Stop() {
	clear(p.Tasks)
	if p.playerState == Stopped {
		return
	}
	if p.playerState == Paused {
		p.Resume()
	}
	speaker.Clear()
	for _, music := range p.MusicList {
		music.Streamer.Close()
	}
	p.CurrentMusicIndex = 0
	p.MusicList = make([]Music, 0)
	p.playerState = Stopped
}

func (p *Player) Pause() {
	if p.playerState == Paused || p.playerState == Stopped {
		return
	}
	p.playerState = Paused
	speaker.Lock()
}

func (p *Player) Resume() {
	if p.playerState == Playing || p.playerState == Stopped {
		return
	}
	p.playerState = Playing
	speaker.Unlock()
}

func (p *Player) Seek(d time.Duration) {
	if p.playerState == Stopped {
		return
	}
	speaker.Lock()
	currentSamplePos := p.MusicList[p.CurrentMusicIndex].Streamer.Position()
	curretnTimePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.D(currentSamplePos)
	newTimePos := (curretnTimePos + d) % p.GetCurrentMusicLength()
	newSamplePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.N(newTimePos)
	if err := p.MusicList[p.CurrentMusicIndex].Streamer.Seek(newSamplePos); err != nil {
		fmt.Println(err)
	}
	speaker.Unlock()
}

func (p *Player) GetCurrentMusic() Music {
	return p.MusicList[p.CurrentMusicIndex]
}

func (p *Player) GetCurrentMusicPosition() time.Duration {
	if p.playerState == Stopped {
		return 0
	}
	currentSamplePos := p.MusicList[p.CurrentMusicIndex].Streamer.Position()
	curretnTimePos := p.MusicList[p.CurrentMusicIndex].Format.SampleRate.D(currentSamplePos)
	return curretnTimePos
}

func (p *Player) GetCurrentMusicLength() time.Duration {
	if len(p.MusicList) == 0 {
		return 0
	}
	music := p.GetCurrentMusic()
	return music.Format.SampleRate.D(music.Streamer.Len())
}

func (p *Player) getMusicList() []string {
	musicList := make([]string, 0)
	for _, music := range p.MusicList {
		musicList = append(musicList, music.Name)
	}
	return musicList
}

func (p *Player) CheckWhatIsThis(unknown string) string {
	// check if its a dir or file
	if fi, err := os.Stat(unknown); err == nil {
		if err == nil {
			if fi.IsDir() {
				// check if there is music files in the dir
				files, err := os.Open(unknown)
				if err != nil {
					return "unknown"
				}

				entries, err := files.Readdir(0)
				if err != nil {
					return "unknown"
				}

				for _, entry := range entries {
					if !entry.IsDir() {
						isMp3, _ := p.Converter.IsMp3(unknown + "/" + entry.Name())
						if isMp3 {
							return "dir"
						}
					}
				}
				return "unknown"
			} else {
				isMp3, _ := p.Converter.IsMp3(unknown)
				if isMp3 {
					return "file"
				}
			}
		}
	}
	// check if its queue index
	if index, err := strconv.Atoi(unknown); err == nil {
		if ok := index < len(p.MusicList); ok {
			return "queue"
		}
	}
	// check if its youtube url
	if strings.Contains(unknown, "youtube.com") || strings.Contains(unknown, "youtu.be") {
		return "youtube"
	}
	if len(unknown) == 11 {
		engines := p.Director.GetEngines()
		engine, ok := engines["youtube"]
		if !ok {
			return "unknown"
		}
		exists, err := engine.Exists(unknown)
		if err != nil {
			return "unknown"
		}
		if exists {
			return "youtube"
		}

	}
	return "unknown"
}

func (p *Player) GetAvailableMusicOptions(query string) []shared.SearchResult {
	//add time out
	p.addTask(query, shared.Search)
	defer p.removeTask(query)
	var musics []shared.SearchResult
	var err error
	searchDone := make(chan bool)
	go func() {
		musics, err = p.Director.Search("youtube", query, 5)
		searchDone <- true
		if err != nil {
			p.errorifyTask(query, err)
			log.Println("Failed to search for", query, ":", err)
		}
	}()

	select {
	case <-time.After(60 * time.Second):
		p.errorifyTask(query, errors.New("Timeout searching for "+query))
		log.Println("Timeout searching for", query)
		break
	case <-searchDone:
		close(searchDone)
		break
	}

	// search in cache
	p.Director.Cached.Fetch()
	files := p.Director.Cached.Search(query)
	for _, f := range files {
		name, _ := shared.ParseCachedFileName(strings.Split(f, "/")[len(strings.Split(f, "/"))-1])
		musics = append(
			musics,
			shared.SearchResult{
				Title:       name,
				Type:        "cache",
				Destination: f,
			},
		)
	}
	// etc ...
	return musics
}

func (p *Player) DetectAndPlay(unknown string) []shared.SearchResult {
	fmt.Println("Checking what is this")
	whatIsThis := p.CheckWhatIsThis(unknown)
	switch whatIsThis {
	case "youtube":
		log.Println("Detected youtube")
		go p.AddMusicFromYoutube(unknown)
	case "dir":
		log.Println("Detected dir")
		go p.AddMusicsFromDir(unknown)
	case "file":
		log.Println("Detected file")
		go p.AddMusicFromFile(unknown)
	case "queue":
		log.Println("Detected queue")
		p.CurrentMusicIndex, _ = strconv.Atoi(unknown)
		p.Play()
	case "unknown":
		log.Println("Detected unknown, searching for", unknown)
		return p.GetAvailableMusicOptions(unknown)
	}
	return []shared.SearchResult{}
}

func (p *Player) GetPlayerStatus() shared.Status {
	return shared.Status{
		CurrentMusicIndex:    p.CurrentMusicIndex,
		CurrentMusicPosition: p.GetCurrentMusicPosition(),
		CurrentMusicLength:   p.GetCurrentMusicLength(),
		PlayerState:          p.playerState,
		MusicList:            p.getMusicList(),
		Tasks:                p.Tasks,
	}
}

func (p *Player) Wait() {
	<-p.done
}
