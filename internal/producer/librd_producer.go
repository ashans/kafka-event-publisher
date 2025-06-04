package producer

import (
	"context"
	"ev_pub/internal/config"
)

type LibrdProducer struct {
}

func (l *LibrdProducer) Init(ctx context.Context, config config.ProducerConfig) error {

	return nil
}

func (l *LibrdProducer) Produce(ctx context.Context, key, value []byte, headers map[string]string, partition int32) (offset int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (l *LibrdProducer) TopicPartitions(ctx context.Context, topic string) (int32, error) {
	//TODO implement me
	panic("implement me")
}
