package configuration

import (
	"fmt"
	cfg "github.com/PiotrFerenc/mash2/internal/consts"
	"github.com/spf13/viper"
)

type configuration struct {
}

func CreateYmlConfiguration() Configuration {
	return &configuration{}
}

func (config *configuration) LoadConfiguration() *Config {
	viper.SetConfigName(cfg.ConfigurationFileName)
	viper.AddConfigPath(cfg.ConfigurationFolderName)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		QueueHost: viper.GetString(cfg.QueueName),
	}
}
