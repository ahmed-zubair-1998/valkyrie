package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Hub struct {
	topics map[int]*Topic
}

func NewHub() *Hub {
	return &Hub{
		topics: make(map[int]*Topic),
	}
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
	go newTopic.Run()
	return newTopic
}

func (hub *Hub) SubscribeToTopic(w http.ResponseWriter, r *http.Request) {
	topicId, err := strconv.Atoi(r.URL.Query().Get("topicId"))
	if err != nil {
		fmt.Fprint(w, "Invalid topic id provided")
		return
	}

	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprint(w, "Websocket upgrade failed")
		conn.Close()
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
	fmt.Fprint(w, "Successfully subscribed to topic id:", topicId)
}

type BroadcastEventData struct {
	Id      int    `json:"id"`
	Message string `json:"msg"`
}

func (hub *Hub) BroadcastEvent(w http.ResponseWriter, r *http.Request) {
	var data BroadcastEventData
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		fmt.Fprint(w, "Invalid post data")
		return
	}

	event, ok := hub.topics[data.Id]
	if !ok {
		fmt.Fprint(w, "Incorrect event id")
		return
	}

	event.broadcast <- []byte(strconv.Itoa(data.Id) + ":" + data.Message)
}
