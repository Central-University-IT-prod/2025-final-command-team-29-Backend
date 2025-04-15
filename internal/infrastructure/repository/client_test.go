package repository_test

import (
	"testing"

	"backend/internal/infrastructure/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       dbMock,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestNewClientRepository_Valid(t *testing.T) {
	db, _ := setupMockDB(t)
	repo := repository.NewClientRepository(db)
	assert.NotNil(t, repo)
}
func TestNewClientRepository_Nil(t *testing.T) {
	assert.Panics(t, func() {
		repository.NewClientRepository(nil)
	}, "Expected panic when db is nil")
}
