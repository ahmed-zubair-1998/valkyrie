package main

type Topic struct {
	id         int
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (topic *Topic) Run() {
	for {
		select {
		case client := <-topic.register:
			topic.clients[client] = true
		case message := <-topic.broadcast:
			for client := range topic.clients {
				client.send <- message
			}
		}
	}
}
