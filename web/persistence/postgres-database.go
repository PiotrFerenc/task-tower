package persistence

import (
	"fmt"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/web/types"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type postgres struct {
	Database *gorm.DB
	config   *configuration.DatabaseConfig
}

func CreatePostgresDatabase(config *configuration.DatabaseConfig) Database {
	return &postgres{
		config: config,
	}
}

func (p *postgres) Connect() *gorm.DB {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.config.DbUser, p.config.DbPassword, p.config.DbHost, p.config.DbPort, p.config.DbName)

	db, err := gorm.Open(driver.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	p.Database = db
	return db
}

func (p *postgres) RunMigration() {
	err := p.Database.AutoMigrate(&types.Pipeline{}, &types.Parameters{}, &types.Step{})
	if err != nil {
		log.Fatalln(err)
	}
}
