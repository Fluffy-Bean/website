package sse

type Connection struct {
	ID        int64
	Name      string
	Messages  chan Message
	Heartbeat chan bool
}

func NewConnection(username string) *Connection {
	return &Connection{
		Name:      username,
		Messages:  make(chan Message),
		Heartbeat: make(chan bool),
	}
}

func (c *Connection) QueueMessage(message Message) {
	c.Messages <- message
}

func (c *Connection) Keepalive() {
	c.Heartbeat <- true
}
