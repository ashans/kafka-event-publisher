package di

import (
	"context"
	"errors"
	"ev_pub/internal/codec"
	"ev_pub/internal/config"
	"ev_pub/internal/producer"
	"log"
	"plugin"
)

type PluginResolvableContainer struct {
	config      config.AppConfig
	producer    producer.Producer
	encoders    map[string]codec.Encoder
	partitioner map[string]producer.Partitioner
}

func NewPluginResolvableContainer(appConfig config.AppConfig, kafkaProducer producer.Producer) *PluginResolvableContainer {
	return &PluginResolvableContainer{
		config:      appConfig,
		producer:    kafkaProducer,
		encoders:    make(map[string]codec.Encoder),
		partitioner: make(map[string]producer.Partitioner),
	}
}

func (p *PluginResolvableContainer) SetEncoder(key string, encoder codec.Encoder) {
	p.encoders[key] = encoder
}

func (p *PluginResolvableContainer) GetEncoder(key string) (codec.Encoder, error) {
	if enc, ok := p.encoders[key]; ok {
		return enc, nil
	}
	return nil, errors.New("encoder not found with key: " + key)
}

func (p *PluginResolvableContainer) SetPartitioner(key string, partitioner producer.Partitioner) {
	p.partitioner[key] = partitioner
}

func (p *PluginResolvableContainer) GetPartitioner(key string) (producer.Partitioner, error) {
	if partitioner, ok := p.partitioner[key]; ok {
		return partitioner, nil
	}
	return nil, errors.New("partitioner not found with key: " + key)
}

func (p *PluginResolvableContainer) DefaultProducer() producer.Producer {
	return p.producer
}

func (p *PluginResolvableContainer) InitModules(ctx context.Context) error {
	err := p.loadPlugins()
	if err != nil {
		return err
	}
	for key, enc := range p.encoders {
		if init, ok := enc.(Initialize); ok {
			err = init.Init(ctx, p.config.Encoders[key].Configs())
			if err != nil {
				return err
			}
			log.Default().Print("encoder " + key + " has been initialized")
		}
	}

	for key, partitioner := range p.partitioner {
		if init, ok := partitioner.(Initialize); ok {
			err = init.Init(ctx, p.config.Partitioners[key].Configs())
			if err != nil {
				return err
			}
			log.Default().Print("partitioner " + key + " has been initialized")
		}
	}

	if init, ok := p.producer.(InitializeProducer); ok {
		err = init.Init(ctx, p.config.Producer)
		if err != nil {
			return err
		}
		log.Default().Print("producer has been initialized")
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

func (p *PluginResolvableContainer) loadPlugins() error {
	for key, pluginLoadData := range p.config.Plugins.Encoders {
		err := p.loadEncoder(key, pluginLoadData)
		if err != nil {
			return err
		}
		log.Default().Print(`encoder plugin ` + key + ` loaded`)
	}

	for key, pluginLoadData := range p.config.Plugins.Partitioners {
		err := p.loadPartitioner(key, pluginLoadData)
		if err != nil {
			return err
		}
		log.Default().Print(`partitioner plugin ` + key + ` loaded`)
	}

	return nil
}

func (p *PluginResolvableContainer) loadEncoder(key string, loadData config.PluginLoadConfig) error {
	sym, err := p.loadPlugin(loadData)
	if err != nil {
		return err
	}

	encoderPlugin, ok := sym.(codec.Encoder)
	if !ok {
		return errors.New("encoder " + key + " is not a valid encoder plugin type")
	}
	p.SetEncoder(key, encoderPlugin)
	return nil
}

func (p *PluginResolvableContainer) loadPartitioner(key string, loadData config.PluginLoadConfig) error {
	sym, err := p.loadPlugin(loadData)
	if err != nil {
		return err
	}

	partitionerPlugin, ok := sym.(producer.Partitioner)
	if !ok {
		return errors.New("partitioner " + key + " is not a valid partitioner plugin type")
	}
	p.SetPartitioner(key, partitionerPlugin)
	return nil
}

func (p *PluginResolvableContainer) loadPlugin(conf config.PluginLoadConfig) (plugin.Symbol, error) {
	loadedPlugin, err := plugin.Open(`./` + p.config.Plugins.Dir + `/` + conf.FileName)
	if err != nil {
		return nil, err
	}
	sym, err := loadedPlugin.Lookup(conf.Symbol)
	if err != nil {
		return nil, err
	}

	return sym, nil
}
