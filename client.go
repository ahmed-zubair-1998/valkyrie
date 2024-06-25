package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	topic *Topic
	conn  WebsocketInterface
	send  chan []byte
}

func (client *Client) ListenToEvents() {
	for message := range client.send {
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}
