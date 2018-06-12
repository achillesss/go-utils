package wsserver

import (
	"fmt"

	gomap "bitbucket.org/magmeng/go-utils/go-map"
	"github.com/achillesss/go-utils/log"
)

type roomatesMap map[int]struct{}

type room struct {
	id             int
	broadcastChan  chan []byte
	msgReceiveChan chan map[int][]byte
	roomates       *gomap.GoMap
}
type roomsMap map[int]*room

func (r room) ListRoomates() []int {
	var ids []int
	rmap := r.roomates.Interface().(roomatesMap)
	for k := range rmap {
		ids = append(ids, k)
	}
	return ids
}

func (server WsServer) ListRooms() []int {
	var ids []int
	rmap := server.rooms.Interface().(roomsMap)
	for k := range rmap {
		ids = append(ids, k)
	}
	return ids
}

func (server *WsServer) queryRoom(roomID int) *room {
	var r *room
	server.rooms.Query(roomID, &r)
	return r
}

func (server *WsServer) QueryRoomates(roomID int) []int {
	r := server.queryRoom(roomID)
	if r == nil {
		return nil
	}
	return r.ListRoomates()
}

func (r *room) addRoomate(id int) {
	r.roomates.Add(id, struct{}{})
}

func (server *WsServer) AddRoomate(roomID, socketID int) error {
	r := server.queryRoom(roomID)
	if r == nil {
		return fmt.Errorf("room not found")
	}

	s := server.querySocket(socketID)
	if s == nil {
		return fmt.Errorf("socket not found")
	}

	r.addRoomate(socketID)
	return nil
}

func (server *WsServer) RemoveRoomate(roomID, socketID int) error {
	r := server.queryRoom(roomID)
	if r == nil {
		return fmt.Errorf("room not found")
	}
	r.removeRoomate(socketID)
	return nil
}

func (r *room) removeRoomate(id int) {
	r.roomates.Delete(id)
}

func (server WsServer) Broadcast(msg []byte, roomID int) error {
	r := server.queryRoom(roomID)
	if r == nil {
		return fmt.Errorf("room not found")
	}
	ids := r.ListRoomates()
	server.BatchSend(msg, ids...)
	return nil
}

func (server WsServer) chatMonitor() {
	for rmsg := range server.roomMsg {
		go func(rmsg map[int][]byte) {
			for id, msg := range rmsg {
				err := server.Broadcast(msg, id)
				if err != nil {
					log.Errorfln("broadcast %s to %d failed. error: %v", msg, id, err)
				}
			}
		}(rmsg)
	}
}

func (server *WsServer) NewRoom() {
	var r room
	*server.lastRoomID++
	r.id = *server.lastRoomID
	r.broadcastChan = make(chan []byte)
	r.msgReceiveChan = make(chan map[int][]byte)
	r.roomates = gomap.NewMap(make(roomatesMap))
	go r.roomates.Handler()
	server.rooms.Add(r.id, &r)
	log.Infofln("add new room %d", r.id)
}
