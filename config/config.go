package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	KafkaHost                 string `envconfig:"KAFKA_HOST" default:"localhost"`
	KafkaPort                 string `envconfig:"KAFKA_PORT" default:"9092"`
	KafkaProductTopic         string `envconfig:"KAFKA_PRODUCT_TOPIC" default:"product"`
	KafkaProductConsumerGroup string `envconfig:"KAFKA_PRODUCT_CONSUMER_GROUP" default:"product-consumer-group"`

	TokopediaURL   string `envconfig:"TOKOPEDIA_URL" default:"https://a1384ac0-a120-47f2-a5c3-f518085745c0.mock.pstmn.io/v3/products/"`
	OmnichannelURL string `envconfig:"OMNICHANNEL_URL" default:"https://localhost:8080//api/v1/product/marketplace/"`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)

	return cfg
}
