package player

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/Malwarize/goplay/shared"
	"github.com/gopxl/beep/mp3"
)

func (p *Player) RPCPlay(_ int, reply *int) error {
	log.Println("RPCPlay called")
	p.Play()
	*reply = 1
	log.Println("RPCPlay done")
	return nil
}

func (p *Player) RPCNext(_ int, reply *int) error {
	log.Println("RPCNext called")
	p.Next()
	*reply = 1
	log.Println("RPCNext done")
	return nil
}

func (p *Player) RPCPrev(_ int, reply *int) error {
	log.Println("RPCPrev called")
	p.Prev()
	*reply = 1
	log.Println("RPCPrev done")
	return nil
}

func (p *Player) RPCPause(_ int, reply *int) error {
	log.Println("RPCPause called")
	p.Pause()
	*reply = 1
	log.Println("RPCPause done")
	return nil
}

func (p *Player) RPCStop(_ int, reply *int) error {
	log.Println("RPCStop called")
	p.Stop()
	*reply = 1
	log.Println("RPCStop done")
	return nil
}

func (p *Player) RPCGetCurrentMusic(_ int, reply *Music) error {
	log.Println("RPCGetCurrentMusic called")
	*reply = p.GetCurrentMusic()
	log.Println("RPCGetCurrentMusic done with reply :", reply)
	return nil
}

func (p *Player) RPCGetCurrentMusicPosition(_ int, reply *int) error {
	log.Println("RPCGetCurrentMusicPosition called")
	*reply = int(p.GetCurrentMusicPosition().Seconds())
	log.Println("RPCGetCurrentMusicPosition done with reply :", reply)
	return nil
}

func (p *Player) RPCGetCurrentMusicLength(_ int, reply *int) error {
	log.Println("RPCGetCurrentMusicLength called")
	*reply = int(p.GetCurrentMusicLength().Seconds())
	log.Println("RPCGetCurrentMusicLength done with reply :", reply)
	return nil
}

func (p *Player) RPCSeek(d time.Duration, reply *int) error {
	log.Println("RPCSeek called with duration in seconds :", d*time.Second)
	p.Seek(d * time.Second)
	*reply = 0
	log.Println("RPCSeek done")
	return nil
}

func (p *Player) RPCResume(_ int, reply *int) error {
	log.Println("RPCResume called")
	p.Resume()
	*reply = 1
	log.Println("RPCResume done")
	return nil
}

func (p *Player) RPCAddMusic(music string, reply *int) error {
	log.Println("RPCAddMusic called with music :", music)
	f, err := os.Open(music)
	if err != nil {
		fmt.Println(err)
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.AddMusic(NewMusic(music, streamer, format))
	*reply = 1
	log.Println("RPCAddMusic done")
	return nil

}

func (p *Player) RPCGetPlayerStatus(_ int, reply *shared.Status) error {
	log.Println("RPCGetPlayerStatus called")
	*reply = p.GetPlayerStatus()
	log.Println("RPCGetPlayerStatus done with reply :", reply)
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
