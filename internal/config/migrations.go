package config

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func LoadMigrations(source string) {
	m, err := migrate.New(source, viper.GetString("DATABASE_URL"))

	if err != nil {
		Logger.Fatalw("DB migrations has failed", "error", err)
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			Logger.Debug("No DB migrations was applied")
			return
		}

		Logger.Fatalw("DB migrations has failed", "error", err)
	}
}
