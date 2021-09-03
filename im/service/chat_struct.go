package service

import (
	"golang.org/x/net/websocket"
)

type (
	ImClient struct {
		ID     string
		Socket *websocket.Conn
		Send   chan []byte
	}

	ImClientManager struct {
		ImClientMap map[string]*ImClient
		Broadcast   chan []byte
		Register    chan *ImClient
		UnRegister  chan *ImClient
	}
)

var ImManager = ImClientManager{
	ImClientMap: make(map[string]*ImClient),
	Broadcast:   make(chan []byte),
	Register:    make(chan *ImClient),
	UnRegister:  make(chan *ImClient),
}

type ImOnlineMsg struct {
	Code        int    `json:"code, omitempty"`
	Msg         string `json:"msg, omitempty"`
	ID          string `json:"id, omitempty"`
	ChannelType int    `json:"channel_type"`
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

type Msg struct {
	Code        int    `json:"code,omitempty"`
	FromId      int    `json:"from_id,omitempty"`
	Msg         string `json:"msg,omitempty"`
	ToId        int    `json:"to_id,omitempty"`
	Status      int    `json:"status,omitempty"`
	MsgType     int    `json:"msg_type,omitempty"`
	ChannelType int    `json:"channel_type"`
}

type ImMessage struct {
	ID          uint64 `json:"id"`
	Msg         string `json:"msg"`
	CreatedAt   string `json:"created_at"`
	FromId      int    `json:"user_id"`
	ToId        int    `json:"send_id"`
	Channel     string `json:"channel"`
	IsRead      int    `json:"is_read"`
	MsgType     int    `json:"msg_type"`
	ChannelType int    `json:"channel_type"`
}

type OnlineMsg struct {
	Code        int    `json:"code,omitempty"`
	Msg         string `json:"msg,omitempty"`
	ID          string `json:"id,omitempty"`
	ChannelType int    `json:"channel_type"`
}

const (
	connOut = 5000
	connOk  = 1000
	SendOk  = 200
)

type GroupId struct {
	UserId string `json:"user_id"`
}

type GroupMap struct {
	GroupIds map[int]*GroupId
}
