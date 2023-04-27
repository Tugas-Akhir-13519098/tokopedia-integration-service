package main

import (
	"log"
    "github.com/segmentio/kafka-go"
    "context"
    "fmt"
)

func consumeMessages() {
    // Set up the Kafka reader with the appropriate configuration
    config := kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "test",
        GroupID: "test-consumer-group",
    }
    reader := kafka.NewReader(config)

    // Continuously read messages from Kafka
    for {
        msg, err := reader.ReadMessage(context.Background())
        if err != nil {
            fmt.Println("Error reading message from Kafka:", err.Error())
            continue
        }
        fmt.Println("Received message from Kafka:", string(msg.Value))
    }
}

func main() {
	log.Printf("Application is running")
    consumeMessages()
}