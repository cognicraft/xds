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
	direction  Direction
	sensors    map[string]config
}

func (s *XDS) Run() error {
	c, err := s.queue.Connect("$xds")
	if err != nil {
		return err
	}
	defer c.Close()

	s.connection = c
	c.Subscribe("sensor/+/manifest", 0, mqtt.HandlerFunc(s.handleSensorManifest))
	c.Subscribe("sensor/+/temperature", 0, mqtt.HandlerFunc(s.handleSensorTemperature))
	go s.simulate()

	return s.queue.ListenAndServe()
}

func (s *XDS) handleSensorManifest(c mqtt.Connection, m mqtt.Message) {
	ts := time.Now().UTC().Format(time.RFC3339)
	sID := m.Topic.Parts()[1]
	fmt.Printf("%s sensor=%s, manifest=%s\n", ts, sID, string(m.Payload))
}

func (s *XDS) handleSensorTemperature(c mqtt.Connection, m mqtt.Message) {
	ts := time.Now().UTC().Format(time.RFC3339)
	sID := m.Topic.Parts()[1]
	fmt.Printf("%s sensor=%s, temperature=%s\n", ts, sID, string(m.Payload))
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
