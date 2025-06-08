package publisher

import (
	"context"
	"ev_pub/internal/common"
	"ev_pub/internal/di"
	"ev_pub/internal/errors"
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
		return 0, 0, errors.Wrap(err, `error getting key encoder`)
	}
	encodedKey, err := keyEncoder.Encode(ctx, key.Options(), []byte(key.Val()))
	if err != nil {
		return 0, 0, errors.Wrap(err, `error encoding key`)
	}

	valEncoder, err := p.container.GetEncoder(value.Type())
	if err != nil {
		return 0, 0, errors.Wrap(err, `error getting value encoder`)
	}
	encodedVal, err := valEncoder.Encode(ctx, value.Options(), []byte(value.Val()))
	if err != nil {
		return 0, 0, errors.Wrap(err, `error encoding value`)
	}

	partitioner, err := p.container.GetPartitioner(partitioning.Type())
	if err != nil {
		return 0, 0, errors.Wrap(err, `error getting partitioner`)
	}
	partitionCount, err := p.container.DefaultProducer().TopicPartitions(ctx, topic)
	if err != nil {
		return 0, 0, errors.Wrap(err, `error getting topic partition count`)
	}
	partition, err = partitioner.Partition(ctx, partitioning.Options(), []byte(key.Val()), []byte(value.Val()),
		encodedKey, encodedVal, headers, partitionCount)
	if err != nil {
		return 0, 0, errors.Wrap(err, `error finding partition to publish`)
	}

	offset, err = p.container.DefaultProducer().Produce(ctx, topic, encodedKey, encodedVal, headers, partition)
	if err != nil {
		return 0, 0, errors.Wrap(err, `error publishing message`)
	}

	return partition, offset, nil
}
