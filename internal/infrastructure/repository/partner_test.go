package repository_test

import (
	"testing"

	"backend/internal/infrastructure/repository"

	"github.com/stretchr/testify/assert"
)

func TestNewPartnerRepository_Valid(t *testing.T) {
	db, _ := setupMockDB(t)
	repo := repository.NewPartnerRepository(db)
	assert.NotNil(t, repo)
}

func TestNewPartnerRepository_Nil(t *testing.T) {
	assert.Panics(t, func() {
		repository.NewPartnerRepository(nil)
	}, "Database connection is nil in repository")
}
