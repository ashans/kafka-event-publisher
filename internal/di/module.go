package di

import (
	"context"
	"ev_pub/internal/config"
)

type Initialize interface {
	Init(ctx context.Context, configs map[string]string) error
}

type InitializeProducer interface {
	Init(ctx context.Context, config config.ProducerConfig) error
}

type Closable interface {
	Close() error
}
