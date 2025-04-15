package repository

import (
	"backend/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePartner(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPartnerRepository(db)

	partner := &models.Partner{
		ID:           uuid.MustParse("8D01ED56-4EB5-493E-B26F-B36C023D30C5"),
		PhoneNumber:  "1234567890",
		PasswordHash: "kksmdfksmdfsdlfksdkmsf",
		Title:        "Test partner",
		Description:  "Test description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "partner"`).
		WithArgs(partner.ID, partner.PhoneNumber, partner.PasswordHash, partner.Title, partner.Description). // Правильный порядок
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CreatePartner(partner)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPartnerByPhoneNumber(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewPartnerRepository(db)

	phoneNumber := "1234567890"

	rows := sqlmock.NewRows([]string{"phone_number", "title"}).
		AddRow(phoneNumber, "Test partner")

	mock.ExpectQuery(`SELECT .* FROM "partner"`).
		WithArgs(phoneNumber, phoneNumber[1:], 1).
		WillReturnRows(rows)

	client, err := repo.GetPartnerByPhoneNumber(phoneNumber)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, phoneNumber, client.PhoneNumber)
	assert.Equal(t, "Test partner", client.Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}
