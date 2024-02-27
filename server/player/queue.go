package player

import (
	"sync"
)

type MusicQueue struct {
	queue   []Music
	current int
	mu      *sync.Mutex
}

func NewMusicQueue() *MusicQueue {
	return &MusicQueue{
		queue:   make([]Music, 0),
		current: 0,
		mu:      &sync.Mutex{},
	}
}

func (q *MusicQueue) GetTitles() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	titles := make([]string, 0)
	for _, music := range q.queue {
		titles = append(titles, music.Name)
	}
	return titles
}

func (q *MusicQueue) GetCurrIndex() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.current
}

func (q *MusicQueue) SetCurrIndex(index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.current = index
}

func (q *MusicQueue) GetMusicByIndex(index int) *Music {
	if index < 0 || index >= q.Size() {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return &q.queue[index]
}

func (q *MusicQueue) GetMusicByName(name string) *Music {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, music := range q.queue {
		if music.Name == name {
			return &music
		}
	}
	return nil
}

func (q *MusicQueue) GetCurrMusic() *Music {
	if q.IsEmpty() {
		return nil
	}
	mu := q.GetMusicByIndex(q.GetCurrIndex())
	return mu
}

func (q *MusicQueue) Enqueue(music Music) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// check if music is already in queue
	for _, m := range q.queue {
		if hash(m.Data) == hash(music.Data) {
			return
		}
	}
	q.queue = append(q.queue, music)
}

func (q *MusicQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue)
}

func (q *MusicQueue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *MusicQueue) Clear() {
	for _, music := range q.queue {
		music.Streamer().Close()
	}
	q.queue = make([]Music, 0)
	q.SetCurrIndex(0)
}

func (q *MusicQueue) Remove(music *Music) {
	// get index of music
	index := -1
	for i, m := range q.queue {
		if m.Name == music.Name {
			index = i
		}
	}
	if index < 0 || index >= q.Size() {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue[index].Streamer().Close()
	q.queue = append(q.queue[:index], q.queue[index+1:]...)
}

func (q *MusicQueue) QueueNext() {
	index := q.GetCurrIndex() + 1
	if index > q.Size()-1 {
		index = 0
	}
	q.SetCurrIndex(index)
}

func (q *MusicQueue) QueuePrev() {
	index := q.GetCurrIndex() - 1
	if index < 0 {
		index = q.Size() - 1
	}
	q.SetCurrIndex(index)
}
