package config

import (
	"database/sql"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	db, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))
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
		log.Default().Fatalf("postgres connection has failed: %v", err)
		return nil
	}

	if err = db.Ping(); err != nil {
		log.Default().Fatalf("ping has failed: %v", err)
		return nil
	}

	return conn
}
