package main

import (
	"context"
	"ev_pub/internal/codec"
	"ev_pub/internal/config"
	"ev_pub/internal/di"
	"ev_pub/internal/http"
	"ev_pub/internal/producer"
	"ev_pub/internal/publisher"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var appConfig config.AppConfig
	err = yaml.Unmarshal(configFile, &appConfig)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	container := di.NewPluginResolvableContainer(appConfig, new(producer.LibrdProducer))

	container.SetEncoder(`byte`, codec.ByteEncoder{})
	container.SetPartitioner(`defined`, producer.DefinedPartitioner{})

	err = container.InitModules(ctx)
	if err != nil {
		panic(err)
	}

	pub := publisher.NewPublisher(container)

	router := http.NewRouter(appConfig, pub)
	err = router.Start()
	if err != nil {
		panic(err)
	}
}
