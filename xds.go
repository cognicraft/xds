package xds

import (
	"fmt"
	"time"

	"github.com/cognicraft/mqtt"
)

func New(queueBind string) *XDS {
	return &XDS{
		queue: mqtt.NewQueue(queueBind),
	}
}

type XDS struct {
	queue      *mqtt.Queue
	connection mqtt.Connection
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
		if pos, err := s.sensorPosition(sID); err == nil {
			fmt.Printf("%s position=%s, temperature=%s\n", ts, pos, string(data))
		}
	default:
		fmt.Printf("%s %s %s\n", ts, topic, string(data))
	}
}

func (s *XDS) sensorPosition(sID string) (string, error) {
	switch sID {
	case "hbd-1":
		return "left-outer-bearing", nil
	case "hbd-2":
		return "right-outer-bearing", nil
	default:
		return "", fmt.Errorf("unknown position: %s", sID)
	}
}
