package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var WebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

func main() {
	hub := NewHub()

	http.HandleFunc("/heartbeat", Heartbeat)
	http.HandleFunc("/topics/subscribe", hub.SubscribeToTopic)
	http.HandleFunc("/events/broadcast", hub.BroadcastEvent)

	http.ListenAndServe(":8080", nil)
}
