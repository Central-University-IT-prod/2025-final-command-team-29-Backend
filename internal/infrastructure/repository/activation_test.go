package repository

import (
	"regexp"
	"testing"
	"time"

	"backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetActivation(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)

	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	req := &models.ActivationRequest{ClientID: clientID, PromoID: promoID}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "activations" WHERE "activations"."client_id" = $1 AND "activations"."promo_id" = $2 AND "activations"."client_id" = $3 AND "activations"."promo_id" = $4 ORDER BY "activations"."client_id" LIMIT $5`)).
		WithArgs(clientID, promoID, clientID, promoID, 1).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "activations" ("client_id","promo_id","created_at","updated_at","buy","current","overall") VALUES ($1,$2,$3,$4,$5,$6,$7)`)).
		WithArgs(clientID, promoID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	act, err := repo.GetActivation(req)
	assert.NoError(t, err)
	assert.Equal(t, clientID, act.ClientID)
	assert.Equal(t, promoID, act.PromoID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckGlobalConditions(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)

	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	req := &models.ActivationRequest{PromoID: promoID}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "promos" WHERE "promos"."id" = $1 ORDER BY "promos"."id" LIMIT $2`)).
		WithArgs(promoID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "time_limit_start", "time_limit_end", "usage_limit"}).
			AddRow(promoID, nil, nil, nil))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "activations" WHERE promo_id = $1`)).
		WithArgs(promoID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	err := repo.CheckGlobalConditions(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUniqueConditions(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)

	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	req := &models.ActivationRequest{ClientID: clientID, PromoID: promoID}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "activations" WHERE client_id = $1 AND promo_id = $2 ORDER BY "activations"."client_id" LIMIT $3`)).
		WithArgs(clientID, promoID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	act, err := repo.CheckUniqueConditions(req)
	assert.NoError(t, err)
	assert.Equal(t, clientID, act.ClientID)
	assert.Equal(t, promoID, act.PromoID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivateCommon(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)
	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	req := &models.ActivationRequest{ClientID: clientID, PromoID: promoID}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "activations" WHERE "activations"."client_id" = $1 AND "activations"."promo_id" = $2 ORDER BY "activations"."client_id" LIMIT $3`)).
		WithArgs(clientID, promoID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"client_id", "promo_id", "overall", "created_at", "updated_at", "buy", "current",
		}).AddRow(clientID, promoID, 0, time.Now(), time.Now(), nil, nil))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"activations\" SET").
		WithArgs(1, clientID, promoID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := repo.ActivateCommon(req)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActivateBundle(t *testing.T) {
	db, _ := setupMockDB(t)
	repo := NewActivationRepository(db)
	act := &models.Activation{Overall: 5}
	repo.ActivateBundle(&models.ActivationRequest{}, act)
	assert.Equal(t, 5, act.Overall)
	assert.NotNil(t, act.Current)
	assert.Equal(t, 1, *act.Current)
}

func TestActivate(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)

	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	act := &models.Activation{ClientID: clientID, PromoID: promoID, Overall: 5}
	mock.ExpectBegin()

	mock.ExpectExec("UPDATE \"activations\" SET").
		WithArgs(
			sqlmock.AnyArg(), // client_id
			sqlmock.AnyArg(), // promo_id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // buy
			sqlmock.AnyArg(), // current
			6,                // overall (5 + 1)
			promoID,          // WHERE promo_id
			clientID,         // WHERE client_id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Activate(act)
	assert.NoError(t, err)
	assert.Equal(t, 6, act.Overall)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewActivationRepository(db)

	clientID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	promoID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	act := &models.Activation{ClientID: clientID, PromoID: promoID, Overall: 0}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"activations\"").
		WithArgs(clientID, promoID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(act)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
