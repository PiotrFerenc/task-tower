package configuration

type QueueConfig struct {
	QueueRunPipe      string
	QueueTasksucceed  string
	QueueTaskFailed   string
	QueueTaskFinished string
	QueueHost         string
	QueuePort         string
	QueueVhost        string
	QueueUser         string
	QueuePassword     string
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
	DbPort     string
	DbName     string
}

type Configuration interface {
	LoadConfiguration() *Config
}
