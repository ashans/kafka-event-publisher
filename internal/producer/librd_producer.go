package producer

import (
	"context"
	"ev_pub/internal/config"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type LibrdProducer struct {
	producer *kafka.Producer
	admin    *kafka.AdminClient
}

func (l *LibrdProducer) Init(_ context.Context, _ config.ProducerConfig) (err error) {
	l.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return err
	}

	l.admin, err = kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *LibrdProducer) Close() error {
	l.producer.Flush(15 * 1000)
	l.producer.Close()

	l.admin.Close()

	return nil
}

func (l *LibrdProducer) Produce(_ context.Context, topic string, key, value []byte, headers map[string]string,
	partition int32) (offset int64, err error) {
	recordHeaders := make([]kafka.Header, len(headers))
	index := 0
	for k, v := range headers {
		recordHeaders[index] = kafka.Header{Key: k, Value: []byte(v)}
		index++
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err = l.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            key,
		Timestamp:      time.Now(),
		Headers:        recordHeaders,
	}, deliveryChan)
	if err != nil {
		return 0, err
	}

	e := <-deliveryChan
	m, ok := e.(*kafka.Message)
	if !ok {
		return 0, fmt.Errorf("Delivery failed, expected message type: %T, got: %T", m, e)
	}
	if m.TopicPartition.Error != nil {
		return 0, m.TopicPartition.Error
	}

	return int64(m.TopicPartition.Offset), nil
}

func (l *LibrdProducer) TopicPartitions(ctx context.Context, topic string) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	describeTopicsResult, err := l.admin.DescribeTopics(ctx, kafka.NewTopicCollectionOfTopicNames([]string{topic}))
	if err != nil {
		return 0, err
	}

	// Print results
	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Error.Code() != 0 {
			return 0, fmt.Errorf("%s", t.Error.Error())
		}

		return int32(len(t.Partitions)), nil
	}

	return 0, fmt.Errorf("topic %s not found", topic)
}
