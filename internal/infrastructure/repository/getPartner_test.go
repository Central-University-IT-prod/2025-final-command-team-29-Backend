package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPartnersList(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := PartnerRepository{db: db}

	page := 1
	limit := 2

	rows := sqlmock.NewRows([]string{"id", "phone_number", "password_hash", "title", "description"}).
		AddRow("8D01ED56-4EB5-493E-B26F-B36C023D30C5", "1234567890", "hash1", "Partner 1", "Description 1").
		AddRow("4E7F5D2B-0E7E-4C89-80A3-1A0A7B9C6F5F", "0987654321", "hash2", "Partner 2", "Description 2")

	mock.ExpectQuery(`SELECT .* FROM "partner"`).
		WillReturnRows(rows)

	partners, err := repo.GetPartnersList(page, limit)
	assert.NoError(t, err)
	assert.Len(t, partners, 2)

	assert.Equal(t, "1234567890", partners[0].PhoneNumber)
	assert.Equal(t, "Partner 1", partners[0].Title)

	assert.Equal(t, "0987654321", partners[1].PhoneNumber)
	assert.Equal(t, "Partner 2", partners[1].Title)

	assert.NoError(t, mock.ExpectationsWereMet())
}
