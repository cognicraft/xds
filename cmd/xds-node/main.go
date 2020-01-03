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
	c.OnMessage(on)
	time.Sleep(time.Millisecond)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var bs []byte
	bs, _ = json.Marshal(map[string]interface{}{"type": "voestalpine:hbd:1", "name": "HBD 1"})
	c.Publish("sensor/hbd-1/info", bs)
	bs, _ = json.Marshal(map[string]interface{}{"type": "voestalpine:hbd:1", "name": "HBD 2"})
	c.Publish("sensor/hbd-2/info", bs)
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
			c.Publish("sensor/hbd-1/temperature", bs)
			bs, _ = json.Marshal(float64(rand.Intn(150) + 20))
			c.Publish("sensor/hbd-2/temperature", bs)
		}
	}
}

func publishTemperature(c mqtt.Connection, sid int, t float64) {
	prefix := fmt.Sprintf("sensor/%d", sid)
	c.Publish(prefix+"/chan/0", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/1", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/2", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/3", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/4", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/5", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/6", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/7", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/8", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/9", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/10", []byte(fmt.Sprintf("%.0f", t)))
	c.Publish(prefix+"/chan/11", []byte(fmt.Sprintf("%.0f", t)))
}

func on(topic string, data []byte) {
	fmt.Printf("on: %s - %s\n", topic, string(data))
}
