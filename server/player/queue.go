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
		titles = append(titles, music.Name())
	}
	return titles
}

func (q *MusicQueue) GetCurrentIndex() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.current
}

func (q *MusicQueue) SetCurrentIndex(index int) {
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

func (q *MusicQueue) GetCurrentMusic() *Music {
	if q.IsEmpty() {
		return nil
	}
	mu := q.GetMusicByIndex(q.GetCurrentIndex())
	return mu
}

func (q *MusicQueue) Enqueue(music Music) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// check if music is already in queue
	for _, m := range q.queue {
		if music.Path == m.Path {
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
	q.SetCurrentIndex(0)
}

func (q *MusicQueue) Remove(index int) {
	if index < 0 || index >= q.Size() {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue[index].Streamer().Close()
	q.queue = append(q.queue[:index], q.queue[index+1:]...)
}

func (q *MusicQueue) QueueNext() {
	index := q.GetCurrentIndex() + 1
	if index > q.Size()-1 {
		index = 0
	}
	q.SetCurrentIndex(index)
}

func (q *MusicQueue) QueuePrev() {
	index := q.GetCurrentIndex() - 1
	if index < 0 {
		index = q.Size() - 1
	}
	q.SetCurrentIndex(index)
}
