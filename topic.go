package main

type Topic struct {
	id         int
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	stop       chan bool
}

func (topic *Topic) Run() {
	for {
		select {
		case client := <-topic.register:
			topic.clients[client] = true
		case message := <-topic.broadcast:
			for client := range topic.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(topic.clients, client)
				}
			}
		case <-topic.stop:
			break
		}
	}
}
