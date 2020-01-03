package xds

import (
	"fmt"
	"time"

	"github.com/cognicraft/mqtt"
)

func New(queueBind string) *XDS {
	return &XDS{
		queue:     mqtt.NewQueue(queueBind),
		direction: Direction1,
		sensors: map[string]config{
			"1.1": {
				id: "1.1",
				d1: "left-outer-bearing",
				d2: "right-outer-bearing",
			},
			"1.2": {
				id: "1.2",
				d1: "right-outer-bearing",
				d2: "left-outer-bearing",
			},
		},
	}
}

type XDS struct {
	queue      *mqtt.Queue
	connection mqtt.Connection

	direction Direction
	sensors   map[string]config
}

func (s *XDS) Run() error {
	c, err := s.queue.Connect("$xds")
	if err != nil {
		return err
	}
	defer c.Close()

	s.connection = c
	c.Subscribe("#", 0)
	c.OnMessage(s.onMessage)

	go s.simulate()

	return s.queue.ListenAndServe()
}

func (s *XDS) onMessage(topic string, data []byte) {
	ts := time.Now().UTC().Format(time.RFC3339)
	t := mqtt.Topic(topic)
	switch {
	case mqtt.Topic("sensor/+/info").Accept(t):
		sID := t.Parts()[1]
		fmt.Printf("%s sensor=%s, info=%s\n", ts, sID, string(data))
	case mqtt.Topic("sensor/+/temperature").Accept(t):
		sID := t.Parts()[1]
		if pos, ok := s.sensorPosition(sID); ok {
			fmt.Printf("%s position=%s, temperature=%s\n", ts, pos, string(data))
		}
	default:
		fmt.Printf("%s %s %s\n", ts, topic, string(data))
	}
}

func (s *XDS) sensorPosition(sID string) (string, bool) {
	c, ok := s.sensors[sID]
	if !ok {
		return "", false
	}
	switch s.direction {
	case Direction1:
		return c.d1, true
	case Direction2:
		return c.d2, true
	default:
		return "", false
	}
}

func (s *XDS) simulate() {
	for {
		select {
		case <-time.After(time.Second * 10):
			switch s.direction {
			case Direction1:
				s.direction = Direction2
			case Direction2:
				s.direction = Direction1
			}
		}
	}
}

type config struct {
	id string
	d1 string
	d2 string
}

type Direction uint8

const (
	Direction1 = 1
	Direction2 = 2
)
