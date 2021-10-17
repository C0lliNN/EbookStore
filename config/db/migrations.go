package db

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func LoadMigrations(source string) {
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USERNAME")
	pass := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DATABASE")
	ssl := viper.GetString("POSTGRES_SSL")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbName, ssl)

	m, err := migrate.New(source, dbUrl)

	if err != nil {
		log.Logger.Fatalw("DB migration has failed", "error", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Logger.Debug("No DB migration was applied")
			return
		}

		log.Logger.Fatalw("DB migration has failed", "error", err)
	}
}
