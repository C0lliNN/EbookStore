package persistence_test

import (
	"context"
	"time"

	"github.com/ebookstore/internal/migrator"
	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
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
	s.db = config.NewConnection()

	require.Eventually(s.T(), func() bool {
		db, err := s.db.DB()
		if err != nil {
			return false
		}

		return db.Ping() == nil
	}, time.Second*10, time.Millisecond*100)

	m := migrator.New(migrator.Config{
		DatabaseURI: migrator.DatabaseURI(s.container.URI),
		Source:      "file:../../migrations",
	})

	m.Sync()

}

func (s *RepositoryTestSuite) TearDownSuite() {
	ctx := context.TODO()

	_ = s.container.Terminate(ctx)
}
