package publisher

import (
	"context"
	"ev_pub/internal/common"
	"ev_pub/internal/di"
)

type Publisher struct {
	container di.PublisherContainer
}

func NewPublisher(container di.PublisherContainer) *Publisher {
	return &Publisher{
		container: container,
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, key, value, partitioning common.TypeValWithOptions,
	headers map[string]string) (partition int32, offset int64, err error) {
	keyEncoder, err := p.container.GetEncoder(key.Type())
	if err != nil {
		return 0, 0, err
	}
	encodedKey, err := keyEncoder.Encode(ctx, key.Options(), []byte(key.Val()))
	if err != nil {
		return 0, 0, err
	}

	valEncoder, err := p.container.GetEncoder(value.Type())
	if err != nil {
		return 0, 0, err
	}
	encodedVal, err := valEncoder.Encode(ctx, value.Options(), []byte(value.Val()))
	if err != nil {
		return 0, 0, err
	}

	partitioner, err := p.container.GetPartitioner(partitioning.Type())
	if err != nil {
		return 0, 0, err
	}
	partitionCount, err := p.container.DefaultProducer().TopicPartitions(ctx, topic)
	if err != nil {
		return 0, 0, err
	}
	partition, err = partitioner.Partition(ctx, partitioning.Options(), []byte(key.Val()), []byte(value.Val()),
		encodedKey, encodedVal, headers, partitionCount)
	if err != nil {
		return 0, 0, err
	}

	offset, err = p.container.DefaultProducer().Produce(ctx, topic, encodedKey, encodedVal, headers, partition)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
