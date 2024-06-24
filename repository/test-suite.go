package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/security-testing-api/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB            *gorm.DB
	mock          sqlmock.Sqlmock
	hasher        helper.Hasher
	repository    UserRepository
	kYCRepository KYCRepository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dsn := "sqlmock_db_0"
	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.hasher = helper.MockHasher{}
	s.repository = NewUserRepositoryImpl(s.DB, s.hasher)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
