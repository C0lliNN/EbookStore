package config

import (
	"github.com/ebookstore/internal/migrator"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func NewMigrationDatabaseURI() migrator.DatabaseURI {
	return migrator.DatabaseURI(viper.GetString("DATABASE_URI"))
}

func NewMigrationSource() migrator.Source {
	return migrator.Source(viper.GetString("MIGRATION_SOURCE"))
}
