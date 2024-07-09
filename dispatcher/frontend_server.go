package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type FrontendServer struct {
	topics    map[int]bool
	conn      *websocket.Conn
	broadcast chan BroadcastEventData
}

type SubscribeTopicData struct {
	TopicId int `json:"topic_id"`
}

func (server *FrontendServer) BroadcastEvents() {
	for data := range server.broadcast {
		server.SendRequest(data)
	}
}

func (server *FrontendServer) SendRequest(data BroadcastEventData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error before dispatching event", err.Error())
		return
	}

	err = server.conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error while dispatching event", err.Error())
		return
	}
}

func (server *FrontendServer) SubscribeToTopics() {
	for {
		_, message, err := server.conn.ReadMessage()
		if websocket.IsCloseError(err) {
			log.Println("Websocket connection closed by FE server", err)
			return
		}
		if err != nil {
			log.Println("Unable to read message from FE server", err)
			return
		}

		var data SubscribeTopicData
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Println("Unable to read message from FE server", err)
			return
		}

		server.topics[data.TopicId] = true
	}
}
