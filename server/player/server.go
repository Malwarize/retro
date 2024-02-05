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
	go func() {
		log.Println("RPCPlay called")
		p.Play()
		*reply = 1
		log.Println("RPCPlay done")
	}()
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

func (p *Player) RPCRemoveMusic(index int, reply *int) error {
	go func() {
		log.Println("RPCRemoveMusic called with index :", index)
		p.Remove(index)
		*reply = 1
		log.Println("RPCRemoveMusic done")
	}()
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

func (p *Player) RPCPlayListsNames(_ int, reply *[]string) error {
	log.Println("RPCPlayLists called")
	*reply = p.PlayListsNames()
	log.Println("RPCPlayLists done with reply :", reply)
	return nil
}

func (p *Player) RPCCreatePlayList(name string, reply *int) error {
	log.Println("RPCCreatePlaylist called with name :", name)
	p.CreatePlayList(name)
	*reply = 1
	log.Println("RPCCreatePlaylist done")
	return nil
}

func (p *Player) RPCRemovePlayList(name string, reply *int) error {
	log.Println("RPCRemovePlaylist called with name :", name)
	p.RemovePlayList(name)
	*reply = 1
	log.Println("RPCRemovePlaylist done")
	return nil
}

func (p *Player) RPCDetectAndAddToPlayList(args shared.AddToPlayListArgs, reply *[]shared.SearchResult) error {
	log.Println("RPCDetectAndAddToPlayList called with query :", args.Query, " and playlist name :", args.PlayListName)
	*reply = p.DetectAndAddToPlayList(args.PlayListName, args.Query)
	log.Println("RPCDetectAndAddToPlayList done")
	return nil
}

func (p *Player) RPCPlayListSongs(name string, reply *[]string) error {
	log.Println("RPCPlayListSongs called with name :", name)
	*reply = p.PlayListSongs(name)
	log.Println("RPCPlayListSongs done with reply :", reply)
	return nil
}

func (p *Player) RPCRemoveSongFromPlayList(args shared.RemoveSongFromPlayListArgs, reply *int) error {
	log.Println("RPCRemoveSongFromPlayList called with name :", args.PlayListName, "and index :", args.Index)
	p.RemoveSongFromPlayList(args.PlayListName, args.Index)
	*reply = 1
	log.Println("RPCRemoveSongFromPlayList done")
	return nil
}

func (p *Player) RPCPlayListPlaySong(args shared.PlayListPlaySongArgs, reply *int) error {
	log.Println("RPCPlayListPlaySong called with name :", args.PlayListName, " and index :", args.Index)
	p.PlayListPlaySong(args.PlayListName, args.Index)
	*reply = 1
	log.Println("RPCPlayListPlaySong done")
	return nil
}

func (p *Player) RPCPlayListPlayAll(name string, reply *int) error {
	log.Println("RPCPlayListPlayAll called with name :", name)
	p.PlayListPlayAll(name)
	*reply = 1
	log.Println("RPCPlayListPlayAll done")
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
