package conductor

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SerialPortName string `split_words:"true" required:"true"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
