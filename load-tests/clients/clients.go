package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/myzhan/boomer"
)

var bindHost string
var bindPort string
var stopChannel chan bool

func receiveMessages(c *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			boomer.RecordFailure("WSR", "Error", 0, err.Error())
			return
		}

		if strings.HasPrefix(string(message), "Successfully subscribed to topic") {
			boomer.RecordSuccess("WSR", "Connected", 0, 0)
		} else {
			parts := strings.Split(string(message), ":")
			start, err := strconv.Atoi(parts[1])
			if err != nil {
				boomer.RecordFailure("WSR", "Error", 0, err.Error())
			}
			now := time.Now().UnixMilli()
			boomer.RecordSuccess("WSR", "Success", now-int64(start), int64(len(message)))
		}
	}
}

func worker() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: bindHost + ":" + bindPort, Path: "/topics/subscribe"}
	params := u.Query()
	params.Add("topicId", "1")
	u.RawQuery = params.Encode()

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		boomer.RecordFailure("WSR", "Connection Failed", 0, err.Error())
		return
	}
	defer c.Close()

	done := make(chan struct{})
	go receiveMessages(c, done)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	task := &boomer.Task{
		Name: "WSR",
		Fn:   worker,
	}

	boomer.Events.Subscribe(boomer.EVENT_SPAWN, func(workers int, spawnRate float64) {
		stopChannel = make(chan bool)
	})

	boomer.Events.Subscribe(boomer.EVENT_STOP, func() {
		close(stopChannel)
	})

	boomer.Events.Subscribe(boomer.EVENT_QUIT, func() {
		close(stopChannel)
		time.Sleep(time.Second)
	})

	boomer.Run(task)
}

func init() {
	flag.StringVar(&bindHost, "host", "localhost", "host")
	flag.StringVar(&bindPort, "port", "8080", "port")
}
