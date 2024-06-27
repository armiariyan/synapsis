package postgresql

import (
	"fmt"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/labstack/gommon/color"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.PostgresqlDB) (db *gorm.DB) {
	// * if somehow use local db without password, remove the string "password=%s" below
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if cfg.Debug {
		db = db.Debug()
	}

	color.Println(color.Green(fmt.Sprintf("⇨ connected to postgresql db on %s\n", cfg.Name)))
	return
}

func CreateUUIDExtension(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}
