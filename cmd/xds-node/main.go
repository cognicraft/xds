package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/cognicraft/mqtt"
)

func main() {
	c, err := mqtt.Dial("client", "127.0.0.1:1883")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var bs []byte
	bs, _ = json.Marshal(map[string]interface{}{"type": "voestalpine:hbd:1", "name": "HBD 1"})
	c.Publish("sensor/1.1/manifest", bs)
	bs, _ = json.Marshal(map[string]interface{}{"type": "voestalpine:hbd:1", "name": "HBD 2"})
	c.Publish("sensor/1.2/manifest", bs)
	for {
		select {
		case s := <-signals:
			switch s {
			case os.Interrupt:
				c.Close()
				return
			}
		case <-time.After(time.Second * 2):
			bs, _ = json.Marshal(float64(rand.Intn(150) + 20))
			c.Publish("sensor/1.1/temperature", bs)
			bs, _ = json.Marshal(float64(rand.Intn(150) + 20))
			c.Publish("sensor/1.2/temperature", bs)
		}
	}
}

func handleMQTT(c mqtt.Connection, m mqtt.Message) {
	fmt.Printf("on: %s - %s\n", m.Topic, string(m.Payload))
}
