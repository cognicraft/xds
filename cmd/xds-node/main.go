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
	"github.com/cognicraft/xds"
)

func main() {
	c, err := mqtt.Dial("client", "127.0.0.1:1883")
	if err != nil {
		log.Fatal(err)
	}
	c.OnClose(func(c mqtt.Connection) {
		os.Exit(0)
	})
	time.Sleep(time.Millisecond)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var bs []byte
	bs, _ = json.Marshal(xds.Manifest{
		ID:   "1.1",
		Type: "voestalpine:hbd:1",
		Name: "HBD 1",
	})
	c.Publish("sensor/1.1/manifest", bs)
	bs, _ = json.Marshal(xds.Manifest{
		ID:   "1.2",
		Type: "voestalpine:hbd:1",
		Name: "HBD 2",
	})
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
