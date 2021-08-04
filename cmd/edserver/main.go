package main

import (
	"edServer/server"
	"edServer/sse"
	"fmt"
	"os"
)

func brokerFunc(broker *sse.Broker, stat chan []byte) {
	for {
		data := <-stat
		broker.Notifier <- data
	}
}

func main() {
	port := os.Args[1:]
	broker := sse.NewServer()
	e := server.RouteStart(broker)
	go brokerFunc(broker, server.Stat)
	e.Start(fmt.Sprintf(":%s", port[0]))
}
