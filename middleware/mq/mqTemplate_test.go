package mq

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
)

const MQURL = "amqp://guest:guest@localhost:5672/"

type ConnectionPool struct {
	pool chan *amqp.Connection
	mu   sync.Mutex
}

func NewConnectionPool(size int) *ConnectionPool {
	p := &ConnectionPool{
		pool: make(chan *amqp.Connection, size),
	}

	for i := 0; i < size; i++ {
		conn, err := amqp.Dial(MQURL)
		if err != nil {
			log.Fatalf("Failed to create connection: %v", err)
		}
		p.pool <- conn
	}

	return p
}

func (p *ConnectionPool) Get() *amqp.Connection {
	p.mu.Lock()
	defer p.mu.Unlock()
	return <-p.pool
}

func (p *ConnectionPool) Put(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool <- conn
}

type ChannelPool struct {
	pool chan *amqp.Channel
	mu   sync.Mutex
}

func NewChannelPool(size int, conn *amqp.Connection) *ChannelPool {
	p := &ChannelPool{
		pool: make(chan *amqp.Channel, size),
	}

	for i := 0; i < size; i++ {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to create channel: %v", err)
		}
		p.pool <- ch
	}

	return p
}

func (p *ChannelPool) Get() *amqp.Channel {
	p.mu.Lock()
	defer p.mu.Unlock()
	return <-p.pool
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool <- ch
}

func main() {
	connPool := NewConnectionPool(5)
	chPool := NewChannelPool(5, connPool.Get())

	// 使用连接池和通道池
	conn := connPool.Get()
	ch := chPool.Get()

	// ... do something with the connection and channel

	connPool.Put(conn)
	chPool.Put(ch)
}
