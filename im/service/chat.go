package service

import (
	"encoding/json"
	"fmt"
	"fork_go_im/im/http/models/user"
	"fork_go_im/pkg/model"
	"fork_go_im/pkg/pool"
	"fork_go_im/pkg/wordsfilter"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var mutexKey sync.Mutex

func (manager *ImClientManager) ImStart() {
	for {
		select {
		case conn := <-ImManager.Register:
			mutexKey.Lock()
			manager.ImClientMap[conn.ID] = &ImClient{
				ID:     conn.ID,
				Socket: conn.Socket,
				Send:   conn.Send,
			}
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

			pool.AntsPool.Submit(func() {
				func() {
					var msgList []ImMessage
					list := model.DB.Where("to_id=? and is_read=?", id, 0).Find(&msgList)
					if list.Error != nil {
						fmt.Println(list.Error)
					}
					for key := range msgList {
						data, _ := json.Marshal(&Msg{
							Code:        SendOk,
							Msg:         msgList[key].Msg,
							FromId:      msgList[key].FromId,
							ToId:        msgList[key].ToId,
							Status:      0,
							MsgType:     msgList[key].MsgType,
							ChannelType: msgList[key].ChannelType,
						})
						conn.Send <- data
					}
				}()
			})
		case conn := <-ImManager.UnRegister:
			if _, ok := manager.ImClientMap[conn.ID]; ok {
				id, _ := strconv.ParseInt(conn.ID, 10, 64)
				user.SetUserStatus(uint64(id), 0)
				jsonMessage, _ := json.Marshal(
					&OnlineMsg{
						Code:        connOut,
						Msg:         "User is offline now" + conn.ID,
						ID:          conn.ID,
						ChannelType: 3,
					},
				)
				manager.ImSend(jsonMessage, conn)
				conn.Socket.Close()
				close(conn.Send)
				delete(manager.ImClientMap, conn.ID)
			}
		case message := <-ImManager.Broadcast:
			data := EnMessage(message)
			msg := new(Msg)
			err := json.Unmarshal([]byte(data.Content), &msg)
			if err != nil {
				fmt.Println(err)
			}
			jsonMessage, _ := json.Marshal(&Msg{
				Code:        SendOk,
				Msg:         msg.Msg,
				FromId:      msg.FromId,
				ToId:        msg.ToId,
				Status:      0,
				MsgType:     msg.MsgType,
				ChannelType: msg.ChannelType,
			})
			if msg.ChannelType == 1 {
				connId := strconv.Itoa(msg.ToId)
				if data, ok := manager.ImClientMap[connId]; ok {
					pool.AntsPool.Submit(func() {
						PutData(msg, 1, msg.ChannelType)
					})
					data.Send <- jsonMessage
				} else {
					pool.AntsPool.Submit(func() {
						PutData(msg, 0, msg.ChannelType)
					})
				}
			} else {
				groups, _ := GetGroupUid(msg.ToId)
				for _, value := range groups {
					if data, ok := manager.ImClientMap[value.UserId]; ok {
						pool.AntsPool.Submit(func() {
							PutGroupData(msg, 1, msg.ChannelType)
							data.Send <- jsonMessage
						})
					}
				}
			}
		}
	}
}

func (manager *ImClientManager) ImSend(message []byte, ignore *ImClient) {
	data, ok := manager.ImClientMap[ignore.ID]
	fmt.Println(ignore.ID)
	if ok {
		data.Send <- message
	}
}

func (c *ImClient) ImWrite() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *ImClient) ImRead() {
	defer func() {
		ImManager.UnRegister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			ImManager.UnRegister <- c
			c.Socket.Close()
			break
		}
		msg := new(Msg)
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println(err)
		}
		if wordsfilter.MsgFilter(msg.Msg) {
			c.Socket.WriteMessage(
				websocket.TextMessage,
				[]byte(`{"code": 401, "data": "Sensitive words are not allowed."}`),
			)
			continue
		} else {
			if msg.ChannelType == 1 {
				data := fmt.Sprintf(
					`{"code": 200, "msg": "%s", "from_id": %v, "to_id": %v, "status": "0", "msg_type": %v, "channel_type": %v}`,
					msg.Msg, msg.FromId, msg.ToId, msg.MsgType, msg.ChannelType,
				)
				c.Socket.WriteMessage(websocket.TextMessage, []byte(data))
			}
		}
		if string(message) == "HeartBeat" {
			c.Socket.WriteMessage(
				websocket.TextMessage,
				[]byte(`{"code": 0, "data": "heatbeat ok"}`),
			)
			continue
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.ID, Content: string(message)})
		ImManager.Broadcast <- jsonMessage
	}
}
