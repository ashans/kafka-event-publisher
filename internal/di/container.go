package di

import (
	"ev_pub/internal/codec"
	"ev_pub/internal/producer"
)

type PublisherContainer interface {
	SetEncoder(key string, encoder codec.Encoder)
	GetEncoder(key string) (codec.Encoder, error)
	SetPartitioner(key string, partitioner producer.Partitioner)
	GetPartitioner(key string) (producer.Partitioner, error)
	DefaultProducer() producer.Producer
}
