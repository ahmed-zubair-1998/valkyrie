package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type HubInterface interface {
	SubscribeToTopic(w http.ResponseWriter, r *http.Request)
	BroadcastEvent(w http.ResponseWriter, r *http.Request)
}

type WebsocketConnectionInterface interface {
	WriteMessage(messageType int, data []byte) error
}

var WebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

func SetupRoutes(hub HubInterface) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/heartbeat", Heartbeat)
	mux.HandleFunc("/topics/subscribe", hub.SubscribeToTopic)
	mux.HandleFunc("/events/broadcast", hub.BroadcastEvent)

	return mux
}

func main() {
	hub := NewHub()
	mux := SetupRoutes(hub)
	http.ListenAndServe(":8080", mux)
}
