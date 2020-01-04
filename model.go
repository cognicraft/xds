package xds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cognicraft/mqtt"
)

func NewModel() *Model {
	return &Model{
		sensors: map[string]Sensor{},
	}
}

type Model struct {
	mu      sync.RWMutex
	sensors map[string]Sensor
}

func (m *Model) Sensor(mf Manifest) Sensor {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := &TemperatureSensor{
		Manifest:    mf,
		Temperature: 0,
	}
	m.sensors[mf.ID] = s

	return s
}

func (m *Model) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "Sensors\n")
	for _, s := range m.sensors {
		fmt.Fprintf(buf, "  %#v\n", s)
	}
	return buf.String()
}

type Sensor interface {
	mqtt.Handler
}

type TemperatureSensor struct {
	Manifest    Manifest
	Temperature float64
}

func (s *TemperatureSensor) HandleMQTT(c mqtt.Connection, m mqtt.Message) {
	t := float64(0)
	err := json.Unmarshal(m.Payload, &t)
	if err != nil {
		return
	}
	s.Temperature = t
}
