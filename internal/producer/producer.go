package producer

import (
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type ObjectProd struct {
	Prod *kafka.Producer
}

type ProdInterface interface {
	SendMessage(topic string, partNum int, msg ProduceMessage) error
	GenerateKeyUUID(value int) []string
}

type ProduceMessage struct {
	Kye       string
	Value     string
	TimeStamp time.Time
	//Header    []kafka.Header
}

// Interface instance kafka
func InstProd(p *kafka.Producer) (ProdInterface, error) {

	if p == nil {
		return nil, errors.New("missed pointer kafka")
	}
	return &ObjectProd{
		Prod: p,
	}, nil
}

// Create producer
func NewProd(servers string) (*kafka.Producer, func(), error) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("fault create producer:{%v}", err)
	}

	fc := func() {
		p.Close()
	}

	return p, fc, nil
}

// Send message
func (p *ObjectProd) SendMessage(topic string, partNum int, msg ProduceMessage) error {

	chStatusTx := make(chan kafka.Event)

	err := p.Prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(partNum),
		},
		Value:     []byte(msg.Value),
		Key:       []byte(msg.Kye),
		Timestamp: msg.TimeStamp,
		//Headers:        []kafka.Header{{Key: "mykey", Value: []byte("myvalue")}},
	}, chStatusTx)
	if err != nil {
		return err
	}

	e := <-chStatusTx
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return fmt.Errorf("fault Tx message: {%v}", ev)
	default:
		return errors.New("unknown error message produce")
	}
}

// Generate UUID keys
func (p *ObjectProd) GenerateKeyUUID(value int) []string {

	keys := make([]string, 0)

	for i := 0; i < value; i++ {
		keys = append(keys, uuid.NewString())
	}

	return keys
}
