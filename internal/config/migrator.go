package config

import (
	"github.com/c0llinn/ebook-store/internal/migrator"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func NewMigrationDatabaseURL() migrator.DatabaseURL {
	return migrator.DatabaseURL(viper.GetString("DATABASE_URL"))
}

func NewMigrationSource() migrator.Source {
	return "file:../migrations"
}
