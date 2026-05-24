package sse

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"git.leggy.dev/Fluffy/Website/internal/broker"
)

type SSE struct {
	Events      *broker.Broker
	Heartbeat   time.Duration
	counter     atomic.Int64
	mut         sync.Mutex
	cancel      context.CancelFunc
	connections map[int64]*Connection
}

func NewSSE(ctx context.Context) *SSE {
	ctx, cancel := context.WithCancel(ctx)
	heartbeat := 5 * time.Second

	s := &SSE{
		Events:      broker.NewBroker(),
		Heartbeat:   heartbeat,
		connections: make(map[int64]*Connection),
		cancel:      cancel,
	}

	go func() {
		time.Sleep(time.Duration(rand.Int63n(int64(heartbeat / 2))))

		ticker := time.NewTicker(heartbeat)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Keepalive()

			case <-ctx.Done():
				return
			}
		}
	}()

	return s
}

func (s *SSE) Subscribe(connection *Connection) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if connection.ID == 0 {
		connection.ID = s.counter.Add(1)
	}

	s.connections[connection.ID] = connection
}

func (s *SSE) Unsubscribe(connection *Connection) {
	s.mut.Lock()
	defer s.mut.Unlock()

	delete(s.connections, connection.ID)
}

func (s *SSE) Keepalive() {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, conn := range s.connections {
		conn.Keepalive()
	}
}

func (s *SSE) Broadcast(message Message) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, conn := range s.connections {
		conn.QueueMessage(message)
	}
}

func (s *SSE) GetConnection(id int64) *Connection {
	s.mut.Lock()
	defer s.mut.Unlock()

	conn, ok := s.connections[id]
	if !ok {
		return nil
	}

	return conn
}

func (s *SSE) GetConnectionsCount() int {
	s.mut.Lock()
	defer s.mut.Unlock()

	return len(s.connections)
}
