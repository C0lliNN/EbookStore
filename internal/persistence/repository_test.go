package persistence_test

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/migrator"
	"github.com/c0llinn/ebook-store/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	container *test.PostgresContainer
}

func (s *RepositoryTestSuite) SetupSuite() {
	ctx := context.TODO()

	var err error
	s.container, err = test.NewPostgresContainer(ctx)
	if err != nil {
		panic(err)
	}

	viper.SetDefault("DATABASE_URI", s.container.URI)

	m := migrator.New(migrator.Config{
		DatabaseURI: migrator.DatabaseURI(s.container.URI),
		Source: "file:../../migrations",
	})

	m.Sync()

	s.db = config.NewConnection()
}

func (s *RepositoryTestSuite) TearDownSuite() {
	ctx := context.TODO()

	s.container.Terminate(ctx)
}
