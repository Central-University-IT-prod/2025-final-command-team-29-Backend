package repository

import (
	"testing"

	"backend/internal/models"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetStat(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewStatRepository(db)

	partnerID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	expected := models.Stat{
		ActivPromos:      5,
		BuyPromos:        120,
		ActiveUsers:      85,
		LimitPromos:      10,
		DailyBuyPromo:    30,
		DailyActiveUsers: 25,
		Retention30d:     40,
		TotalClients:     500,
		ConversionRate:   17.0,
		TopClients:       10,
	}
	rows := sqlmock.NewRows([]string{
		"activ_promos", "buy_promos", "active_users", "limit_promos",
		"daily_buy_promo", "daily_active_users", "retention30d",
		"total_clients", "conversion_rate", "top_clients",
	}).AddRow(
		expected.ActivPromos, expected.BuyPromos, expected.ActiveUsers, expected.LimitPromos,
		expected.DailyBuyPromo, expected.DailyActiveUsers, expected.Retention30d,
		expected.TotalClients, expected.ConversionRate, expected.TopClients,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT \n    COUNT(DISTINCT CASE WHEN a.current > 0 THEN p.id END) AS activ_promos,")).
		WithArgs(partnerID).
		WillReturnRows(rows)

	stat, err := repo.GetStat(partnerID)
	assert.NoError(t, err)
	assert.Equal(t, expected, stat)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStatByPromo(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewStatRepository(db)

	partnerID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	promoID := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	expected := models.StatPromo{
		BuyPromos:     50,
		Users:         40,
		DailyBuyPromo: 20,
		DailyUsers:    15,
	}
	rows := sqlmock.NewRows([]string{
		"buy_promos", "users", "daily_buy_promo", "daily_users",
	}).AddRow(
		expected.BuyPromos, expected.Users, expected.DailyBuyPromo, expected.DailyUsers,
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT \n            COALESCE(SUM(a.buy), 0) AS buy_promos_count,")).
		WithArgs(promoID, partnerID).
		WillReturnRows(rows)

	statPromo, err := repo.GetStatByPromo(partnerID, promoID)
	assert.NoError(t, err)
	assert.Equal(t, expected, statPromo)
	assert.NoError(t, mock.ExpectationsWereMet())
}
