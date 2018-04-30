package connection

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nats-io/go-nats"

	"github.com/venushome/hands.noolite/engine"
)

// NatsConnector - коннектор к очереди сообщений
type NatsConnector struct {
	URL       string
	System    string
	Subsystem string

	EngineOutput chan engine.Event
	EngineInput  chan engine.Command

	conn     *nats.Conn
	stopChan chan bool
	done     sync.WaitGroup
}

// NewNatsConnector - создание коннектора
func NewNatsConnector(serverURL string) *NatsConnector {
	c := NatsConnector{}
	c.URL = serverURL
	return &c
}

// Run - запускает подключение к очереди сообщений
func (c *NatsConnector) Run() {
	c.stopChan = make(chan bool)
	conn, err := nats.Connect(c.URL)
	if err != nil {
		panic(err)
	}
	c.conn = conn

	topic := "venushome." + c.System + "." + c.Subsystem
	c.conn.Subscribe(topic, c.receiveMessage)
	c.done.Add(1)
	go c.handle()
}

// Stop - останавливает обработчик коммуникации
func (c *NatsConnector) Stop() {
	close(c.stopChan)
	c.done.Wait()
}

func (c *NatsConnector) handle() {
	defer c.done.Done()

	for {
		select {
		case _, ok := <-c.stopChan:
			if !ok {
				conn := c.conn
				c.conn = nil
				conn.Close()
				return
			}
		case evnt := <-c.EngineOutput:
			c.publishMessage(evnt)
		}
	}
}

func (c *NatsConnector) receiveMessage(m *nats.Msg) {
	fmt.Printf("Nats Message %s\n", m.Data)
	cmd := natsCommand{}
	if err := json.Unmarshal(m.Data, &cmd); err != nil {
		return
	}
	fmt.Println("!!!!!!!", cmd)
	if cmd.System == c.System && cmd.Subsystem == c.Subsystem {
		c.EngineInput <- cmd.Payload
	}
}

func (c *NatsConnector) publishMessage(evnt engine.Event) {
	conn := c.conn
	if conn == nil {
		return
	}

	msg := natsEvent{
		System:    c.System,
		Subsystem: c.Subsystem,
		Payload:   evnt,
	}
	bmsg, err := json.Marshal(msg)
	if err != nil {
		return
	}
	fmt.Printf("Nats event %s\n", string(bmsg))
	conn.Publish("venushome.head", bmsg)
}

type natsCommand struct {
	System    string
	Subsystem string
	Payload   engine.Command
}

type natsEvent struct {
	System    string
	Subsystem string
	Payload   engine.Event
}
