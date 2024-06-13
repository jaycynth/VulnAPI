package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestUserRepositoryImpl_Save(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	user := model.User{
		Id:       1,
		Username: "testuser",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Username).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo.Save(user)
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	user := model.User{
		Id:       1,
		Username: "updateduser",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").
		WithArgs(user.Username, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo.Update(user)
}

func TestUserRepositoryImpl_Delete(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	userId := 1

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `users` WHERE `users`.`id` = ?").
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo.Delete(userId)
}

func TestUserRepositoryImpl_FindById(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	user := model.User{
		Id:       1,
		Username: "testuser",
	}

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(user.Id, user.Username)

	mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`id` = ?").
		WithArgs(user.Id).
		WillReturnRows(rows)

	result, err := repo.FindById(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUserRepositoryImpl_FindAll(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	users := []model.User{
		{Id: 1, Username: "user1"},
		{Id: 2, Username: "user2"},
	}

	rows := sqlmock.NewRows([]string{"id", "username"})
	for _, user := range users {
		rows.AddRow(user.Id, user.Username)
	}

	mock.ExpectQuery("SELECT * FROM `users`").
		WillReturnRows(rows)

	result, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, users, result)
}

func TestUserRepositoryImpl_Login(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	user := model.User{
		Id:       1,
		Username: "testuser",
		Password: "password",
	}

	mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`username` = ? AND `users`.`password` = ?").
		WithArgs(user.Username, user.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(user.Id, user.Username, user.Password))

	result, err := repo.Login(user.Username, user.Password)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUserRepositoryImpl_Login_InvalidData(t *testing.T) {
	db, _ := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	_, err := repo.Login("", "")
	assert.Error(t, err)
	assert.Equal(t, helper.ErrInvalidData, err)
}

func TestUserRepositoryImpl_Login_UserNotFound(t *testing.T) {
	db, mock := setupDB(t)
	repo := NewUserRepositoryImpl(db)

	mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`username` = ? AND `users`.`password` = ?").
		WithArgs("nonexistent", "password").
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.Login("nonexistent", "password")
	assert.Error(t, err)
	assert.Equal(t, helper.ErrUserNotFound, err)
	assert.Equal(t, model.User{}, user)
}
