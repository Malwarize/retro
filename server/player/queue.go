package player

type MusicQueue struct {
	queue   []Music
	current int
}

func NewMusicQueue() *MusicQueue {
	return &MusicQueue{
		queue:   make([]Music, 0),
		current: 0,
	}
}

func (q *MusicQueue) GetTitles() []string {
	titles := make([]string, 0)
	for _, music := range q.queue {
		titles = append(titles, music.Name())
	}
	return titles
}

func (q *MusicQueue) GetCurrentIndex() int {
	return q.current
}

func (q *MusicQueue) SetCurrentIndex(index int) {
	q.current = index
}

func (q *MusicQueue) GetCurrentMusic() *Music {
	if q.IsEmpty() {
		return nil
	}
	return &q.queue[q.current]
}

func (q *MusicQueue) Enqueue(music Music) {
	q.queue = append(q.queue, music)
}

func (q *MusicQueue) Size() int {
	return len(q.queue)
}

func (q *MusicQueue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *MusicQueue) Clear() {
	for _, music := range q.queue {
		music.Streamer.Close()
	}
	q.queue = make([]Music, 0)
	q.current = 0
}

func (p *Player) AddMusicToQueue(music Music) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Queue.Enqueue(music)
}

func (p *Player) RemoveMusicFromQueue(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if index < 0 || index >= p.Queue.Size() {
		return
	}
	p.Queue.queue = append(p.Queue.queue[:index], p.Queue.queue[index+1:]...)
}

func (p *Player) GetMusicQueue() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Queue.GetTitles()
}

func (p *Player) GetCurrentMusic() *Music {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Queue.GetCurrentMusic()
}

func (p *Player) QueueNext() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Queue.current++
	if p.Queue.current >= p.Queue.Size() {
		p.Queue.current = 0
	}
}

func (p *Player) QueuePrev() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Queue.current--
	if p.Queue.current < 0 {
		p.Queue.current = p.Queue.Size() - 1
	}
}

func (p *Player) SetCurrentIndex(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Queue.SetCurrentIndex(index)
}
