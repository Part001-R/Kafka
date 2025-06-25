package main

import (
	"fmt"
	"kafka/internal/producer"
	"log"
	"time"
)

func main() {

	bootstrapServers := "localhost:9092,localhost:9093,localhost:9094"

	// Create kafka producer
	p, fc, err := producer.NewProd(bootstrapServers)
	if err != nil {
		log.Fatalf("fault create kafka producer:{%v}", err)
	}
	defer fc()

	instProd, err := producer.InstProd(p)
	if err != nil {
		log.Fatalf("fault create producer instance:{%v}", err)
	}

	// Generation UUID keys
	keysMsg := instProd.GenerateKeyUUID(3)

	// Produce messages to Kafka (Partition - 0)
	for i := 0; i < 60; i++ {
		time.Sleep(100 * time.Millisecond)

		message := fmt.Sprintf("message 1 number %d", i)

		var msg = producer.ProduceMessage{
			Kye:       keysMsg[0],
			Value:     message,
			TimeStamp: time.Now(),
		}
		topic := "test"
		err = instProd.SendMessage(topic, 0, msg)
		if err != nil {
			fmt.Printf("Tx faulr: -> {%s}: error{%v}\n", message, err)
			continue
		}
		fmt.Printf("Tx Ok: -> {%s}\n", message)
	}

	// Produce messages to Kafka (Partition - 1)
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)

		message := fmt.Sprintf("message 2 number %d", i)

		var msg = producer.ProduceMessage{
			Kye:       keysMsg[1],
			Value:     message,
			TimeStamp: time.Now(),
		}

		topic := "test"
		err = instProd.SendMessage(topic, 1, msg)
		if err != nil {
			fmt.Printf("Tx faulr: -> {%s}: error{%v}\n", message, err)
			continue
		}
		fmt.Printf("Tx Ok: -> {%s}\n", message)
	}

	// Produce messages to Kafka (Partition - 2)
	for i := 0; i < 15; i++ {
		time.Sleep(100 * time.Millisecond)

		message := fmt.Sprintf("message 3 number %d", i)

		var msg = producer.ProduceMessage{
			Kye:       keysMsg[2],
			Value:     message,
			TimeStamp: time.Now(),
		}

		topic := "test"
		err = instProd.SendMessage(topic, 2, msg)
		if err != nil {
			fmt.Printf("Tx faulr: -> {%s}: error{%v}\n", message, err)
			continue
		}
		fmt.Printf("Tx Ok: -> {%s}\n", message)
	}

	// Wait for all messages to be delivered
	fmt.Println("wait")
	p.Flush(3000)
}
