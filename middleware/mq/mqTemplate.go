package mq

import (
	"github.com/streadway/amqp"
	"simple_douyin/config"
	"sync"
)

const poolSize = 5

type ConnectionPool struct {
	mu        sync.Mutex
	available []*amqp.Connection
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		available: make([]*amqp.Connection, 0, poolSize),
	}
}

func (p *ConnectionPool) Get() (*amqp.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.available) == 0 {
		conn, err := amqp.Dial(config.MqUrl)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	conn := p.available[0]
	p.available = p.available[1:]
	return conn, nil
}

func (p *ConnectionPool) Put(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.available) >= poolSize {
		conn.Close()
		return
	}

	p.available = append(p.available, conn)
}

type ChannelPool struct {
	mu        sync.Mutex
	available []*amqp.Channel
	conn      *amqp.Connection
}

func NewChannelPool(conn *amqp.Connection) *ChannelPool {
	return &ChannelPool{
		available: make([]*amqp.Channel, 0, poolSize),
		conn:      conn,
	}
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.available) == 0 {
		ch, err := p.conn.Channel()
		if err != nil {
			return nil, err
		}
		return ch, nil
	}

	ch := p.available[0]
	p.available = p.available[1:]
	return ch, nil
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.available) >= poolSize {
		ch.Close()
		return
	}

	p.available = append(p.available, ch)
}
