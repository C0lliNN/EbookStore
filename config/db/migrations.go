package db

import (
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func LoadMigrations(source string) {
	m, err := migrate.New(source, viper.GetString("DATABASE_URL"))

	if err != nil {
		log.Logger.Fatalw("DB migration has failed", "error", err)
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Logger.Debug("No DB migration was applied")
			return
		}

		log.Logger.Fatalw("DB migration has failed", "error", err)
	}
}
