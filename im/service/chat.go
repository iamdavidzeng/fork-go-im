package service

import (
	"encoding/json"
	"fork_go_im/im/http/models/user"
	"strconv"
	"sync"
)

var mutexKey sync.Mutex

func (manager *ImClientManager) ImStart() {
	for {
		select {
		case conn := <-ImManager.Register:
			mutexKey.Lock()
			manager.ImClientMap[conn.ID] = &ImClient{ID: conn.ID, Socket: conn.Socket, Send: conn.Send}
			mutexKey.Unlock()
			jsonMessage, _ := json.Marshal(&ImOnlineMsg{
				Code:        connOk,
				Msg:         "Users are online now.",
				ID:          conn.ID,
				ChannelType: 3,
			})
			id, _ := strconv.ParseInt(conn.ID, 10, 64)
			user.SetUserStatus(uint64(id), 1)
			manager.ImSend(jsonMessage, conn)
		}
	}
}
