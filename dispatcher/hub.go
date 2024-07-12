package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	servers map[WebsocketConnectionInterface]*FrontendServer
}

type BroadcastEventData struct {
	TopicId int    `json:"topic_id"`
	Message string `json:"message"`
}

func NewHub() *Hub {
	return &Hub{
		servers: make(map[WebsocketConnectionInterface]*FrontendServer),
	}
}

func (hub *Hub) GetOrCreateServer(conn WebsocketConnectionInterface) *FrontendServer {
	server, isPresent := hub.servers[conn]
	if isPresent {
		return server
	}

	newServer := &FrontendServer{
		conn:      conn,
		topics:    make(map[int]bool),
		broadcast: make(chan *BroadcastEventData),
	}
	hub.servers[conn] = newServer
	go newServer.BroadcastEvents()
	return newServer
}

func (hub *Hub) FrontendServerConnection(w http.ResponseWriter, r *http.Request) {
	websocketUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, "Websocket upgrade failed")
		return
	}

	server := hub.GetOrCreateServer(conn)
	go server.SubscribeToTopics()
}

func (hub *Hub) BroadcastEvent(w http.ResponseWriter, r *http.Request) {
	var data BroadcastEventData
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		fmt.Fprint(w, "Invalid request data")
		return
	}

	isValidTopic := false
	for _, server := range hub.servers {
		_, found := server.topics[data.TopicId]
		if found {
			server.broadcast <- &data
			isValidTopic = true
		}
	}

	if !isValidTopic {
		fmt.Fprint(w, "Incorrect topic id")
	} else {
		fmt.Fprint(w, "Event broadcasted successfully")
	}
}
