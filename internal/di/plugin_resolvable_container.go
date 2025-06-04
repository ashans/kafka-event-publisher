package di

import (
	"context"
	"ev_pub/internal/codec"
	"ev_pub/internal/config"
	"ev_pub/internal/producer"
)

type PluginResolvableContainer struct {
	config      config.AppConfig
	producer    producer.Producer
	encoders    map[string]codec.Encoder
	partitioner map[string]producer.Partitioner
}

func NewPluginResolvableContainer(appConfig config.AppConfig, producer producer.Producer) *PluginResolvableContainer {
	return &PluginResolvableContainer{
		producer: producer,
		config:   appConfig,
	}
}

func (p *PluginResolvableContainer) SetEncoder(key string, encoder codec.Encoder) {
	p.encoders[key] = encoder
}

func (p *PluginResolvableContainer) GetEncoder(key string) (codec.Encoder, error) {
	return p.encoders[key], nil
}

func (p *PluginResolvableContainer) SetPartitioner(key string, partitioner producer.Partitioner) {
	p.partitioner[key] = partitioner
}

func (p *PluginResolvableContainer) GetPartitioner(key string) (producer.Partitioner, error) {
	return p.partitioner[key], nil
}

func (p *PluginResolvableContainer) DefaultProducer() producer.Producer {
	return p.producer
}

func (p *PluginResolvableContainer) InitModules(ctx context.Context) error {
	for key, enc := range p.encoders {
		if init, ok := enc.(Initialize); ok {
			err := init.Init(ctx, p.config.Encoders[key].Configs())
			if err != nil {
				return err
			}
		}
	}

	for key, partitioner := range p.partitioner {
		if init, ok := partitioner.(Initialize); ok {
			err := init.Init(ctx, p.config.Partitioners[key].Configs())
			if err != nil {
				return err
			}
		}
	}

	if init, ok := p.producer.(InitializeProducer); ok {
		err := init.Init(ctx, p.config.Producer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PluginResolvableContainer) CloseModules(ctx context.Context) error {
	for key, enc := range p.encoders {
		if init, ok := enc.(Initialize); ok {
			err := init.Init(ctx, p.config.Encoders[key].Configs())
			if err != nil {
				return err
			}
		}
	}

	if closable, ok := p.producer.(Closable); ok {
		err := closable.Close()
		if err != nil {
			return err
		}
	}

	for _, partitioner := range p.partitioner {
		if closable, ok := partitioner.(Closable); ok {
			err := closable.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
