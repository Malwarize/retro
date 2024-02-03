package player

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/Malwarize/goplay/shared"
)

func (p *Player) RPCPlay(_ int, reply *int) error {
	log.Println("RPCPlay called")
	p.Play()
	*reply = 1
	log.Println("RPCPlay done")
	return nil
}

func (p *Player) RPCNext(_ int, reply *int) error {
	go func() {
		log.Println("RPCNext called")
		p.Next()
		*reply = 1
		log.Println("RPCNext done")
	}()
	return nil
}

func (p *Player) RPCPrev(_ int, reply *int) error {
	go func() {
		log.Println("RPCPrev called")
		p.Prev()
		*reply = 1
		log.Println("RPCPrev done")
	}()
	return nil
}

func (p *Player) RPCPause(_ int, reply *int) error {
	go func() {
		log.Println("RPCPause called")
		p.Pause()
		*reply = 1
		log.Println("RPCPause done")
	}()
	return nil
}

func (p *Player) RPCStop(_ int, reply *int) error {
	go func() {
		log.Println("RPCStop called")
		p.Stop()
		*reply = 1
		log.Println("RPCStop done")
	}()
	return nil
}

func (p *Player) RPCGetCurrentMusic(_ int, reply *Music) error {
	go func() {
		log.Println("RPCGetCurrentMusic called")
		*reply = p.GetCurrentMusic()
		log.Println("RPCGetCurrentMusic done with reply :", reply)
	}()

	return nil
}

func (p *Player) RPCGetCurrentMusicPosition(_ int, reply *int) error {
	go func() {
		log.Println("RPCGetCurrentMusicPosition called")
		*reply = int(p.GetCurrentMusicPosition().Seconds())
		log.Println("RPCGetCurrentMusicPosition done with reply :", reply)
	}()

	return nil
}

func (p *Player) RPCSeek(d time.Duration, reply *int) error {
	go func() {
		log.Println("RPCSeek called with duration in seconds :", d*time.Second)
		p.Seek(d * time.Second)
		*reply = 0
		log.Println("RPCSeek done")
	}()
	return nil
}

func (p *Player) RPCResume(_ int, reply *int) error {
	go func() {
		log.Println("RPCResume called")
		p.Resume()
		*reply = 1
		log.Println("RPCResume done")
	}()
	return nil
}

func (p *Player) RPCAddMusic(music string, reply *int) error {
	log.Println("RPCAddMusic called with music :", music)
	p.AddMusicFromFile(music)
	*reply = 1
	log.Println("RPCAddMusic done")
	return nil

}

func (p *Player) RPCAddDirectory(path string, reply *int) error {
	log.Println("AddDirectory called with path :", path)
	p.AddMusicsFromDir(path)
	log.Println("AddDirectory done")
	return nil
}

func (p *Player) RPCGetPlayerStatus(_ int, reply *shared.Status) error {
	log.Println("RPCGetPlayerStatus called")
	*reply = p.GetPlayerStatus()
	log.Println("RPCGetPlayerStatus done with reply :", reply)
	return nil
}

func (p *Player) RPCDetectAndPlay(query string, reply *[]shared.SearchResult) error {
	log.Println("RPCDetectAndPlay called with query :", query)
	*reply = p.DetectAndPlay(query)
	log.Println("RPCDetectAndPlay done with reply :", reply)
	return nil
}

func StartIPCServer(port string) {
	log.Println("Creating Player instance")
	player := NewPlayer()
	rpc.Register(player)
	log.Println("Player instance created and registered to RPC")
	lis, err := net.Listen("tcp", ":"+port)

	log.Println("Starting IPC server on ", lis.Addr().String())
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
		}
		go rpc.ServeConn(conn)
	}
}
