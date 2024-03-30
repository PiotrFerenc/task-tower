package configuration

type QueueConfig struct {
	QueueRunPipe       string
	QueueStageSucceed  string
	QueueStageFailed   string
	QueueStageFinished string
	QueueHost          string
	QueuePort          string
	QueueVhost         string
	QueueUser          string
	QueuePassword      string
}
type FolderConfig struct {
	TmpFolder string
}

type Config struct {
	Queue    QueueConfig
	Folder   FolderConfig
	Database DatabaseConfig
}
type DatabaseConfig struct {
	DbHost     string
	DbUser     string
	DbPassword string
}

type Configuration interface {
	LoadConfiguration() *Config
}
