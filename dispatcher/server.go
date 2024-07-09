package main

import (
	"fmt"
	"net/http"
)

type HubInterface interface {
	FrontendServerConnection(w http.ResponseWriter, r *http.Request)
	BroadcastEvent(w http.ResponseWriter, r *http.Request)
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

func SetupRoutes(hub HubInterface) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/heartbeat", Heartbeat)
	mux.HandleFunc("/frontend/connect", hub.FrontendServerConnection)
	mux.HandleFunc("/events/broadcast", hub.BroadcastEvent)

	return mux
}

func main() {
	hub := NewHub()
	mux := SetupRoutes(hub)
	http.ListenAndServe(":8090", mux)
}
