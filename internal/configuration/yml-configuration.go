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

// LoadConfiguration loads the configuration from the specified file and returns a Config object.
// The method uses Viper to read the configuration file.
// It sets the config name and config path using the constants ConfigurationFileName
// and ConfigurationFolderName respectively.
// If the configuration file cannot be read, it panics with a fatal error message.
// The returned Config object contains the following fields:
// - Queue: a QueueConfig object containing various queue-related configurations.
//   - QueueRunPipe: the name of the queue run pipe.
//   - QueueTaskSucceed: the name of the queue task for successful tasks.
//   - QueueTaskFailed: the name of the queue task for failed tasks.
//   - QueueHost: the host of the queue.
//   - QueuePort: the port of the queue.
//   - QueueVhost: the virtual host of the queue.
//   - QueueUser: the username for accessing the queue.
//   - QueuePassword: the password for accessing the queue.
//   - QueueTaskFinished: the name of the queue task for finished tasks.
//
// - Folder: a FolderConfig object containing folder configurations.
//   - TmpFolder: the path of the temporary folder.
//
// - Database: a DatabaseConfig object containing database configurations.
//   - DbHost: the host of the database.
//   - DbUser: the username for accessing the database.
//   - DbPassword: the password for accessing the database.
//   - DbPort: the port of the database.
//   - DbName: the name of the database.
func (config *configuration) LoadConfiguration() *Config {
	viper.SetConfigName(cfg.ConfigurationFileName)
	viper.AddConfigPath(cfg.ConfigurationFolderName)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		Queue: QueueConfig{
			QueueRunPipe:      viper.GetString(cfg.QueueRunPipe),
			QueueTaskSucceed:  viper.GetString(cfg.QueueTasksucceed),
			QueueTaskFailed:   viper.GetString(cfg.QueueTaskFailed),
			QueueHost:         viper.GetString(cfg.QueueHost),
			QueueVhost:        viper.GetString(cfg.QueueVhost),
			QueueUser:         viper.GetString(cfg.QueueUser),
			QueuePassword:     viper.GetString(cfg.QueuePassword),
			QueuePort:         viper.GetString(cfg.QueuePort),
			QueueTaskFinished: viper.GetString(cfg.QueueFinished),
		},
		Folder: FolderConfig{
			TmpFolder: viper.GetString(cfg.TmpFolder),
		},
		Database: DatabaseConfig{
			DbHost:     viper.GetString(cfg.DbHost),
			DbUser:     viper.GetString(cfg.DbUser),
			DbPassword: viper.GetString(cfg.DbPassword),
			DbPort:     viper.GetString(cfg.DbPort),
			DbName:     viper.GetString(cfg.DbName),
		},
		Controller: ControllerConfig{
			Host: viper.GetString(cfg.ControllerHost),
			Post: viper.GetString(cfg.ControllerPort),
		},
	}
}
