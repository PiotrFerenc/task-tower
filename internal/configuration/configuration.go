package configuration

type Config struct {
	QueueHost string
}

type Configuration interface {
	LoadConfiguration() *Config
}
