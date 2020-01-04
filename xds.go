package xds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cognicraft/mqtt"
)

func New(queueBind string) *XDS {
	return &XDS{
		queue: mqtt.NewQueue(queueBind),
		model: NewModel(),
	}
}

type XDS struct {
	queue      *mqtt.Queue
	connection mqtt.Connection

	model *Model
}

func (s *XDS) Run() error {
	c, err := s.queue.Connect("$xds")
	if err != nil {
		return err
	}
	defer c.Close()

	s.connection = c
	c.Subscribe("sensor/+/manifest", 0, mqtt.HandlerFunc(s.handleSensorManifest))
	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				fmt.Printf("\n%s\n", s.model)
			}
		}
	}()

	return s.queue.ListenAndServe()
}

func (s *XDS) handleSensorManifest(c mqtt.Connection, m mqtt.Message) {
	mf := Manifest{}
	err := json.Unmarshal(m.Payload, &mf)
	if err != nil {
		return
	}

	sensor := s.model.Sensor(mf)
	st := mqtt.Topic("sensor/" + mf.ID + "/#")
	go c.Subscribe(st, 0, sensor)
}

type Direction uint8

const (
	Direction1 = 1
	Direction2 = 2
)
