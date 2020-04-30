package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var mqttBroker *Broker

func main() {
	var err error
	mqttBroker, err = NewBroker()
	if err != nil {
		fmt.Printf("MQTT broker init failed! %v\n", err)
	}

	run()
}

// Start the MQTT plugin
func run() {

	go func() {
		if err := mqttBroker.Start(); err != nil {
			fmt.Printf("Stopping MQTT Broker: %s\n", err.Error())
		} else {
			fmt.Printf("Starting MQTT Broker (port %s) ... done\n", mqttBroker.config.Port)
		}
	}()

	if mqttBroker.config.Port != "" {
		fmt.Printf("You can now listen to MQTT via: http://%s:%s\n", mqttBroker.config.Host, mqttBroker.config.Port)
	}

	if mqttBroker.config.TlsPort != "" {
		fmt.Printf("You can now listen to MQTT via: https://%s:%s\n", mqttBroker.config.TlsHost, mqttBroker.config.TlsPort)
	}

	connectClient(mqttBroker.config.Port)

	shutdownSignal := make(chan os.Signal)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal

	fmt.Printf("Stopping MQTT Broker ...\n")

	if err := mqttBroker.Shutdown(); err != nil {
		fmt.Printf("Stopping MQTT Broker: %s\n", err.Error())
	} else {
		fmt.Printf("Stopping MQTT Broker ... done\n")
	}
}
