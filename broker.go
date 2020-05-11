package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/fhmq/hmq/broker"
)

// Broker ... Simple mqtt publisher abstraction
type Broker struct {
	broker *broker.Broker
	config *broker.Config
}

// NewBroker ... Create a new publisher.
func NewBroker() (*Broker, error) {

	wsPort := os.Getenv("PORT")
	if wsPort == "" {
		wsPort = "1880"
	}

	configFilePath := "config.json"
	params := []string{
		fmt.Sprintf("--config=" + configFilePath),
	}

	c, err := broker.ConfigureConfig(params)
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	fmt.Printf("CNOFIG %+v\n", c)

	b, err := broker.NewBroker(c)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}

	return &Broker{
		broker: b,
		config: c,
	}, nil
}

// Start the broker
func (b *Broker) Start() error {
	b.broker.Start()
	return nil
}

// Shutdown the broker.
func (b *Broker) Shutdown() error {
	//return b.broker.Close()
	return nil
}

// Send a msg.
func (b *Broker) Send(topic string, message string) error {

	packet := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	packet.TopicName = topic
	packet.Qos = 0
	packet.Payload = []byte(message)

	b.broker.PublishMessage(packet)

	return nil
}
