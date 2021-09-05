package db

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() (conn *gorm.DB) {
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USERNAME")
	pass := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DATABASE")

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)

	dialector := postgres.New(postgres.Config{
		DSN:                  dbUrl,
		PreferSimpleProtocol: true,
	})
	conn, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Logger.Fatalw("Postgres connection has failed", "error", err.Error())
		return
	}

	db, err := conn.DB()
	if err != nil {
		log.Logger.Fatalw("Postgres connection has failed", "error", err.Error())
		return
	}

	if err = db.Ping(); err != nil {
		log.Logger.Fatalw("Ping has failed", "error", err.Error())
		return
	}

	return conn
}
