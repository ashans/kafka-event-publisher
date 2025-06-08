package common

import (
	"ev_pub/internal/errors"
	"gopkg.in/yaml.v3"
)

type ModuleConfig struct {
	configs map[string]string
}

func (m *ModuleConfig) UnmarshalYAML(value *yaml.Node) error {
	configs := make(map[string]string)
	err := value.Decode(&configs)
	if err != nil {
		return errors.Wrap(err, `error decoding module config`)
	}
	m.configs = configs

	return nil
}

func (m ModuleConfig) Configs() map[string]string {
	return m.configs
}
