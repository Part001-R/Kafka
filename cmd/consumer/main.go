package main

import (
	"fmt"
	"kafka/internal/consumer"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	bootstrapServers := "localhost:9092,localhost:9093,localhost:9094"
	groupID := "test_consumer_1"
	topics := []string{"test"}

	// Create a new Kafka consumer instance
	c, err := consumer.InstanceConsumer(bootstrapServers, groupID)
	if err != nil {
		log.Fatalf("fault create consumer:{%v}", err)
	}
	iCons, err := consumer.InterfaceConsumer(c)
	if err != nil {
		log.Fatalf("fault create consumer interface:{%v}", err)
	}

	// Subscribe partition
	err = iCons.SubscribePartitionNumb(topics[0], 1, kafka.OffsetStored)
	if err != nil {
		fmt.Printf("Error subscribing to specific partition: %v\n", err)
		os.Exit(1)
	}

	// Set up a channel for handling Ctrl-C / SIGINT
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				err := iCons.ShowMessageHandler(e)
				if err != nil {
					log.Fatalf("fault read Kafka message")
				}
				_, err = c.CommitMessage(e)
				if err != nil {
					log.Fatalf("Error committing offset: %v\n", err)
				}

			case kafka.Error:
				log.Fatalf("Error: %v: %v\n", e.Code(), e)

			default:
				fmt.Printf("event: {%v}", ev)
			}
		}
	}

	fmt.Println("Closing consumer")
	c.Close()
}
