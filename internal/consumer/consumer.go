package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
)

type ObjectCons struct {
	Cons *kafka.Consumer
}

type iConsumer interface {
	SubscrTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ShowMessageHandler(e kafka.Event) error
	SubscribePartitionNumb(topic string, partitionNum int32, offset kafka.Offset) error
}

func InterfaceConsumer(c *kafka.Consumer) (iConsumer, error) {

	if c == nil {
		return nil, errors.New("pointer consumer is missed")
	}

	return &ObjectCons{
		Cons: c,
	}, nil
}

func InstanceConsumer(bootstrapServers, groupID string) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, fmt.Errorf("fault create consumer:{%v}", err)
	}

	return c, nil
}

func (c *ObjectCons) SubscribePartitionNumb(topic string, partitionNum int32, offset kafka.Offset) error {
	topicPartition := kafka.TopicPartition{
		Topic:     &topic,
		Partition: partitionNum,
		Offset:    offset,
	}

	err := c.Cons.Assign([]kafka.TopicPartition{topicPartition})
	if err != nil {
		return fmt.Errorf("failed to assign topic partition: %w", err)
	}
	fmt.Printf("Subscribed to topic '%s', partition %d from offset %v\n", topic, partitionNum, offset)
	return nil
}

func (c *ObjectCons) SubscrTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {

	err := c.Cons.SubscribeTopics(topics, rebalanceCb)
	if err != nil {
		return fmt.Errorf("fault subscribe topics:{%v}", err)
	}
	return nil
}

func (c *ObjectCons) ShowMessageHandler(e kafka.Event) error {

	switch ev := e.(type) {
	case *kafka.Message:
		fmt.Printf("Message on %s:\n", ev.TopicPartition.String())
		fmt.Printf("  Timestamp: %s\n", ev.Timestamp.String())
		fmt.Printf("  Key: %s\n", string(ev.Key))
		fmt.Printf("  Value: %s\n", string(ev.Value))
		fmt.Println()
		return nil
	default:
		return fmt.Errorf("fault handler consumer - ricieved not message")
	}
}
