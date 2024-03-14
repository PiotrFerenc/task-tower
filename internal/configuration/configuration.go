package configuration

type QueueConfig struct {
	QueueName     string
	QueueHost     string
	QueuePort     string
	QueueVhost    string
	QueueUser     string
	QueuePassword string
}

type Config struct {
	Queue QueueConfig
}

type Configuration interface {
	LoadConfiguration() *Config
}
