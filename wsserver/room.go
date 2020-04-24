package wsserver

// type roomatesMap map[int64]struct{}
//
// type room struct {
// 	id             int64
// 	broadcastChan  chan []byte
// 	msgReceiveChan chan map[int64][]byte
// 	roomates       *gomap.GoMap
// }
// type roomsMap map[int64]*room
//
// func (r room) ListRoomates() []int64 {
// 	var ids []int64
// 	rmap := r.roomates.Interface().(roomatesMap)
// 	for k := range rmap {
// 		ids = append(ids, k)
// 	}
// 	return ids
// }
//
// func (server WsServer) ListRooms() []int64 {
// 	var ids []int64
// 	rmap := server.rooms.Interface().(roomsMap)
// 	for k := range rmap {
// 		ids = append(ids, k)
// 	}
// 	return ids
// }
//
// func (server *WsServer) queryRoom(roomID int64) *room {
// 	var r *room
// 	server.rooms.Query(roomID, &r)
// 	return r
// }
//
// func (server *WsServer) QueryRoomates(roomID int64) []int64 {
// 	r := server.queryRoom(roomID)
// 	if r == nil {
// 		return nil
// 	}
// 	return r.ListRoomates()
// }
//
// func (r *room) addRoomates(ids ...int64) {
// 	var roomates = map[int64]struct{}{}
// 	for _, id := range ids {
// 		roomates[id] = struct{}{}
// 	}
// 	r.roomates.BatchAdd(roomates)
// }
//
// func (server *WsServer) AddRoomates(roomID int64, socketIDs ...int64) error {
// 	r := server.queryRoom(roomID)
// 	if r == nil {
// 		return fmt.Errorf("room not found")
// 	}
//
// 	socks := server.querySockets(socketIDs...)
// 	var ids []int64
// 	for _, s := range socks {
// 		ids = append(ids, s.id)
// 	}
//
// 	r.addRoomates(ids...)
// 	return nil
// }
//
// func (server *WsServer) RemoveRoomate(roomID, socketID int64) error {
// 	r := server.queryRoom(roomID)
// 	if r == nil {
// 		return fmt.Errorf("room not found")
// 	}
// 	r.removeRoomate(socketID)
// 	return nil
// }
//
// func (r *room) removeRoomate(id int64) {
// 	r.roomates.Delete(id)
// }
//
// func (server *WsServer) Broadcast(msg []byte, roomID int64) error {
// 	r := server.queryRoom(roomID)
// 	if r == nil {
// 		return fmt.Errorf("room not found")
// 	}
// 	ids := r.ListRoomates()
// 	server.SendTo(msg, ids...)
// 	return nil
// }
//
// func (server *WsServer) chatMonitor() {
// 	for rmsg := range server.roomMsg {
// 		go func(rmsg map[int64][]byte) {
// 			for id, msg := range rmsg {
// 				err := server.Broadcast(msg, id)
// 				if err != nil {
// 					log.Errorfln("broadcast %s to room %d failed. error: %v", msg, id, err)
// 				}
// 			}
// 		}(rmsg)
// 	}
// }
//
// func (server *WsServer) NewRoom() int64 {
// 	var r room
// 	r.id = atomic.AddInt64(&server.lastRoomID, 1)
// 	r.broadcastChan = make(chan []byte)
// 	r.msgReceiveChan = make(chan map[int64][]byte)
// 	r.roomates = gomap.NewMap(make(roomatesMap))
// 	server.rooms.Add(r.id, &r)
// 	return r.id
// }
//
// func (server *WsServer) RemoveRoom(roomID int) {
// 	server.rooms.Delete(roomID)
// }
//
// func (server *WsServer) CountRoomates(roomID int) (int, bool) {
// 	var r *room
// 	server.rooms.Query(roomID, &r)
// 	if r == nil {
// 		return 0, false
// 	}
// 	count := r.roomates.Len()
// 	return count, true
// }
