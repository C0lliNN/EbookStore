package config

import (
	"database/sql"

	"github.com/ebookstore/internal/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	db, err := sql.Open("postgres", viper.GetString("DATABASE_URI"))
	if err != nil {
		log.Default().Fatalf("postgres connection has failed: %v", err)
		return nil
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Default().Fatalf("postgres connection has failed: %T %v", err, err)
		return nil
	}

	return conn
}
