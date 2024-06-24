package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	hasher     helper.Hasher
	repository UserRepository
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

func (s *Suite) Test_repository_Save() {
	user := &model.User{
		Username: "johndoe",
		Password: "password123",
		Email:    "johndoe@example.com",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`username`,`password`,`email`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)")).
		WithArgs(user.Username, "$2a$10$e0MYzXyjpJS7Pd0RVvHwHeFUpH0lHLjZ0f72mMBp/Rs8C3W.rFYTO", user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	savedUser, err := s.repository.Save(user)
	require.NoError(s.T(), err)

	require.Equal(s.T(), user.Username, savedUser.Username)
	require.Equal(s.T(), user.Email, savedUser.Email)

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func (s *Suite) Test_repository_Login() {
	username := "johndoe"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username: "johndoe",
		Password: string(hashedPassword),
		Email:    "johndoe@example.com",
	}

	queryPattern := `SELECT \* FROM ` + "`users`" + ` WHERE username = \? ORDER BY ` + "`users`.`id`" + ` LIMIT \?`
	s.mock.ExpectQuery(queryPattern).
		WithArgs(username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"username", "password", "email"}).
			AddRow(user.Username, user.Password, user.Email))

	foundUser, err := s.repository.Login(username, password)
	if err != nil {
		switch {
		case errors.Is(err, helper.ErrMissingData):
			s.T().Errorf("Missing data error: %v", err)
		case errors.Is(err, helper.ErrUserNotFound):
			s.T().Errorf("User not found error: %v", err)
		case errors.Is(err, helper.ErrInvalidData):
			s.T().Errorf("Invalid data error: %v", err)
		default:
			s.T().Errorf("Unexpected error: %v", err)
		}
	} else {
		require.Equal(s.T(), user.Username, foundUser.Username)
		require.Equal(s.T(), user.Email, foundUser.Email)
	}

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_repository_Login_UserNotFound() {
	username := "nonexistentuser"
	password := "password123"

	queryPattern := `SELECT \* FROM ` + "`users`" + ` WHERE username = \? ORDER BY ` + "`users`.`id`" + ` LIMIT \?`
	s.mock.ExpectQuery(queryPattern).WithArgs(username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := s.repository.Login(username, password)
	require.ErrorIs(s.T(), err, helper.ErrUserNotFound)

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_repository_Login_WrongPassword() {
	username := "johndoe"
	password := "wrongpassword"
	hashedPassword, _ := s.hasher.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    "johndoe@example.com",
	}

	queryPattern := `SELECT \* FROM ` + "`users`" + ` WHERE username = \? ORDER BY ` + "`users`.`id`" + ` LIMIT \?`
	s.mock.ExpectQuery(queryPattern).WithArgs(username, 1).
		WillReturnRows(sqlmock.NewRows([]string{"username", "password", "email"}).
			AddRow(user.Username, user.Password, user.Email))

	_, err := s.repository.Login(username, password)
	require.ErrorIs(s.T(), err, helper.ErrInvalidData)

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) Test_repository_FindAll() {
	users := []*model.User{
		{Username: "user1", Password: "password1", Email: "user1@example.com"},
		{Username: "user2", Password: "password2", Email: "user2@example.com"},
	}

	rows := sqlmock.NewRows([]string{"username", "password", "email"})
	for _, user := range users {
		rows.AddRow(user.Username, user.Password, user.Email)
	}

	s.mock.ExpectQuery(`SELECT \* FROM ` + "`users`").
		WillReturnRows(rows)

	foundUsers, err := s.repository.FindAll()
	require.NoError(s.T(), err)
	require.Len(s.T(), foundUsers, len(users))

	for i, user := range foundUsers {
		require.Equal(s.T(), users[i].Username, user.Username)
		require.Equal(s.T(), users[i].Email, user.Email)
	}

	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
