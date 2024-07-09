package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type HubInterface interface {
	SubscribeToTopic(w http.ResponseWriter, r *http.Request)
}

type WebsocketConnectionInterface interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "heartbeat")
}

func SetupRoutes(hub HubInterface) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/heartbeat", Heartbeat)
	mux.HandleFunc("/topics/subscribe", hub.SubscribeToTopic)

	return mux
}

func SetupConnectionToDispatcher(serverAddress string) (*websocket.Conn, error) {
	baseUrl := strings.Replace(serverAddress, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(baseUrl+"/frontend/connect", nil)
	return conn, err
}

func main() {
	conn, err := SetupConnectionToDispatcher("http://localhost:8090")
	if err != nil {
		log.Fatal("Unable to connect to dispatcher", err)
	}
	mux := SetupRoutes(NewHub(conn))
	http.ListenAndServe(":8080", mux)
}
