package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/internal/infrastructure/repository"
)

func TestNewStatRepository_Valid(t *testing.T) {
	db, _ := setupMockDB(t)
	statRepo := repository.NewStatRepository(db)
	assert.NotNil(t, statRepo)
}

func TestNewStatRepository_Nil(t *testing.T) {
	assert.Panics(t, func() {
		repository.NewStatRepository(nil)
	}, "Expected panic when db is nil")
}
