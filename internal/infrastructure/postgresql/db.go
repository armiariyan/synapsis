package postgresql

import (
	"fmt"

	"github.com/armiariyan/synapsis/internal/config"
	"github.com/labstack/gommon/color"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.PostgresqlDB) (db *gorm.DB) {

	//TODO armia clean this
	dsn2 := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	dsn := "host=localhost user=armia password= dbname=online_store sslmode=disable port=5432 sslmode=disable"

	fmt.Println()
	fmt.Println("dsn2")
	fmt.Println(dsn2)
	fmt.Println()
	fmt.Println("dsn")
	fmt.Println(dsn)
	fmt.Println()
	fmt.Println()

	db, err := gorm.Open(postgres.Open(dsn2), &gorm.Config{})
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
