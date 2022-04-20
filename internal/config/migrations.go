package config

import (
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func LoadMigrations(source string) {
	m, err := migrate.New(source, viper.GetString("DATABASE_URL"))

	if err != nil {
		log.Default().Fatalf("migrations has failed: %v", err)
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Default().Debug("No DB migrations was applied")
			return
		}

		log.Default().Fatalf("migrations has failed: %v", err)
	}
}
