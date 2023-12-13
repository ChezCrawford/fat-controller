package conductor

import (
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SerialPortName string `split_words:"true" required:"true"`
	UseSimDriver   bool   `split_words:"true" default:"false"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return c
}
