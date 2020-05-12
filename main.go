package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var mqttBroker *Broker

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println("no .env file")
	}

	mqttBroker, err = NewBroker()
	if err != nil {
		fmt.Printf("MQTT broker init failed! %v\n", err)
	}

	initDB()

	run()
}

// Start the MQTT plugin
func run() {

	router := NewRouter()

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

	// shut down MQTT broker
	if err := mqttBroker.Shutdown(); err != nil {
		fmt.Printf("Stopping MQTT Broker: %s\n", err.Error())
	} else {
		fmt.Printf("Stopping MQTT Broker ... done\n")
	}

	// shutdown web server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting down server: %s", err.Error())
	}
}
