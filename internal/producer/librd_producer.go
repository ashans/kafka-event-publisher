package producer

import (
	"context"
	"ev_pub/internal/config"
	"ev_pub/internal/errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strings"
	"time"
)

type LibrdProducer struct {
	producer *kafka.Producer
	admin    *kafka.AdminClient
}

func (l *LibrdProducer) Init(_ context.Context, config config.ProducerConfig) (err error) {
	bootstrapServers := strings.Join(config.BootstrapServers, `,`)
	l.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return errors.Wrap(err, `failed to create kafka producer`)
	}

	l.admin, err = kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return errors.Wrap(err, `error creating kafka admin client`)
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
		return 0, errors.Wrap(err, `failed to produce message`)
	}

	e := <-deliveryChan
	m, ok := e.(*kafka.Message)
	if !ok {
		return 0, errors.Wrap(fmt.Errorf("delivery failed, expected message type: %T, got: %T", m, e), `failed to produce message`)
	}
	if m.TopicPartition.Error != nil {
		return 0, errors.New(m.TopicPartition.Error.Error())
	}

	return int64(m.TopicPartition.Offset), nil
}

func (l *LibrdProducer) TopicPartitions(ctx context.Context, topic string) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	describeTopicsResult, err := l.admin.DescribeTopics(ctx, kafka.NewTopicCollectionOfTopicNames([]string{topic}))
	if err != nil {
		return 0, errors.Wrap(err, `failed to describe topics`)
	}

	// Print results
	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Error.Code() != 0 {
			return 0, errors.New(t.Error.Error())
		}

		return int32(len(t.Partitions)), nil
	}

	return 0, errors.New("topic " + topic + " not found")
}
