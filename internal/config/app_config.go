package config

import (
	"ev_pub/internal/common"
)

type AppConfig struct {
	AppInfo      map[string]interface{}         `yaml:"app"`
	Http         HttpConfig                     `yaml:"http"`
	Plugins      PluginConfig                   `yaml:"plugins"`
	Producer     ProducerConfig                 `yaml:"producer"`
	Encoders     map[string]common.ModuleConfig `yaml:"encoders"`
	Partitioners map[string]common.ModuleConfig `yaml:"partitioners"`
}

type HttpConfig struct {
	Port int `yaml:"port"`
	Path struct {
		Api string `yaml:"api"`
		Ui  string `yaml:"ui"`
	} `yaml:"path"`
}

type ProducerConfig struct {
	BootstrapServers []string `yaml:"bootstrap_servers"`
}

type PluginConfig struct {
	Dir          string                      `yaml:"dir"`
	Encoders     map[string]PluginLoadConfig `yaml:"encoders"`
	Partitioners map[string]PluginLoadConfig `yaml:"partitioners"`
}

type PluginLoadConfig struct {
	FileName string `yaml:"filename"`
	Symbol   string `yaml:"symbol"`
}
