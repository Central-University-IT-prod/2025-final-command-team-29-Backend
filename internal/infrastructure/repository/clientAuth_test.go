package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/models"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestCreateClient(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewClientRepository(db)

	client := &models.Client{
		ID:           uuid.MustParse("8D01ED56-4EB5-493E-B26F-B36C023D30C5"),
		PhoneNumber:  "1234567890",
		Name:         "Test User",
		PasswordHash: "kksmdfksmdfsdlfksdkmsf",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "client"`).
		WithArgs(client.ID, client.PhoneNumber, client.PasswordHash, client.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreateClient(client)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetClientByPhoneNumber(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewClientRepository(db)

	phoneNumber := "1234567890"

	rows := sqlmock.NewRows([]string{"phone_number", "name"}).
		AddRow(phoneNumber, "Test User")

	mock.ExpectQuery(`SELECT .* FROM "client"`).
		WithArgs(phoneNumber, phoneNumber[1:], 1).
		WillReturnRows(rows)

	client, err := repo.GetClientByPhoneNumber(phoneNumber)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, phoneNumber, client.PhoneNumber)
	assert.Equal(t, "Test User", client.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}
