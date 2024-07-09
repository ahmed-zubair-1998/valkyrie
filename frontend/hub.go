package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Hub struct {
	topics map[int]*Topic
	conn   WebsocketConnectionInterface
}

type BroadcastEventData struct {
	Id      int    `json:"topic_id"`
	Message string `json:"message"`
}

func NewHub(conn WebsocketConnectionInterface) *Hub {
	hub := &Hub{
		topics: make(map[int]*Topic),
		conn:   conn,
	}
	go hub.HandleMessagesFromDispatcher()

	return hub
}

func (hub *Hub) HandleMessagesFromDispatcher() {
	for {
		_, message, err := hub.conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Println("Websocket connection closed by dispatcher", err)
			return
		}
		if err != nil {
			log.Println("Unable to read message from dispatcher", err)
			return
		}

		var data BroadcastEventData
		err = json.Unmarshal(message, &data)
		if err != nil {
			log.Println("Unable to read message from dispatcher", err)
			return
		}

		err = hub.BroadcastEvent(data)
		if err != nil {
			log.Println("Unable to broadcast event", err)
		}
	}
}

func (hub *Hub) BroadcastEvent(data BroadcastEventData) error {
	topic, ok := hub.topics[data.Id]
	if !ok {
		return errors.New("incorrect topic id")
	}

	topic.broadcast <- []byte(strconv.Itoa(data.Id) + ":" + data.Message)
	return nil
}

func (hub *Hub) SubscribeToTopic(w http.ResponseWriter, r *http.Request) {
	topicId, err := strconv.Atoi(r.URL.Query().Get("topicId"))
	if err != nil {
		fmt.Fprint(w, "Invalid topic id provided")
		return
	}

	websocketUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, "Websocket upgrade failed")
		return
	}

	log.Println("Subscribing client to topic id:", topicId)

	topic := hub.GetOrCreateTopic(topicId)
	client := &Client{
		topic: topic,
		conn:  conn,
		send:  make(chan []byte, 256),
	}
	topic.register <- client
	go client.ListenToEvents()

	successMsg := "Successfully subscribed to topic id: " + strconv.Itoa(topicId)
	conn.WriteMessage(websocket.TextMessage, []byte(successMsg))
}

func (hub *Hub) GetOrCreateTopic(id int) *Topic {
	topic, isPresent := hub.topics[id]
	if isPresent {
		return topic
	}

	newTopic := &Topic{
		id:         id,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	hub.topics[id] = newTopic
	go hub.NotifyDispatcher(id)
	go newTopic.Run()
	return newTopic
}

func (hub *Hub) NotifyDispatcher(topicId int) {
	data := struct {
		TopicId int `json:"topic_id"`
	}{
		topicId,
	}
	jsonData, _ := json.Marshal(data)

	err := hub.conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error while notifying dispatcher", err)
	}
}
