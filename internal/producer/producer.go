package producer

import "context"

type Producer interface {
	Produce(ctx context.Context, topic string, key, value []byte, headers map[string]string, partition int32) (offset int64, err error)
	TopicPartitions(ctx context.Context, topic string) (int32, error)
}
