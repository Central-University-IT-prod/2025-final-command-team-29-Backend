package repository

import (
	"fmt"

	"github.com/google/uuid"

	"backend/internal/models"
)

func (r *StatRepository) GetStat(partnerID uuid.UUID) (models.Stat, error) {
	var stat models.Stat

	query := `
       SELECT 
    COUNT(DISTINCT CASE WHEN a.current > 0 THEN p.id END) AS activ_promos,
    COALESCE(SUM(a.buy), 0) AS buy_promos,
    COUNT(DISTINCT a.client_id) AS active_users,
    COUNT(DISTINCT CASE WHEN p.usage_limit IS NOT NULL AND a.buy >= p.usage_limit THEN p.id END) AS limit_promos,
    COALESCE(SUM(CASE WHEN a.buy > 0 AND a.client_id IN (
        SELECT client_id FROM activations WHERE CURRENT_DATE - INTERVAL '1 day' <= CURRENT_DATE
    ) THEN a.buy ELSE 0 END), 0) AS daily_buy_promo,
    COUNT(DISTINCT CASE WHEN a.client_id IN (
        SELECT client_id FROM activations WHERE CURRENT_DATE - INTERVAL '1 day' <= CURRENT_DATE
    ) THEN a.client_id END) AS daily_active_users,
    COUNT(DISTINCT CASE WHEN a.client_id IN (
        SELECT client_id FROM activations 
        WHERE created_at >= NOW() - INTERVAL '30 days'
        GROUP BY client_id HAVING COUNT(*) > 1
    ) THEN a.client_id END) AS retention_30d,
    (SELECT COUNT(*) FROM client) AS total_clients,
    (COUNT(DISTINCT a.client_id) * 100.0 / NULLIF((SELECT COUNT(*) FROM client), 0)) AS conversion_rate,
    COUNT(DISTINCT CASE WHEN a.client_id IN (
        SELECT client_id FROM activations 
        GROUP BY client_id ORDER BY SUM(buy) DESC LIMIT 10
    ) THEN a.client_id END) AS top_clients
FROM promos p
LEFT JOIN activations a ON p.id = a.promo_id
WHERE p.partner_id = ?;
    `

	err := r.db.Raw(query, partnerID).Scan(&stat).Error
	if err != nil {
		return models.Stat{}, err
	}

	return stat, nil
}

func (r *StatRepository) GetStatByPromo(partnerID uuid.UUID, promoID uuid.UUID) (models.StatPromo, error) {
	var stat models.StatPromo
	query := `
        SELECT 
            COALESCE(SUM(a.buy), 0) AS buy_promos_count,
            COUNT(DISTINCT a.client_id) AS users_count,
            COALESCE(SUM(CASE WHEN a.buy > 0 AND a.client_id IN (
                SELECT client_id FROM activations 
                WHERE created_at >= NOW() - INTERVAL '1 day'
            ) THEN a.buy ELSE 0 END), 0) AS daily_buy_promos_count,
            COUNT(DISTINCT CASE WHEN a.client_id IN (
                SELECT client_id FROM activations 
                WHERE created_at >= NOW() - INTERVAL '1 day'
            ) THEN a.client_id END) AS daily_users_count
        FROM activations a
        LEFT JOIN promos p ON a.promo_id = p.id
        WHERE a.promo_id = ? AND p.partner_id = ?
    `
	err := r.db.Raw(query, promoID, partnerID).Scan(&stat).Error
	if err != nil {
		return models.StatPromo{}, err
	}

	return stat, nil
}
func (r *StatRepository) ExportPromoByPartner(partnerID uuid.UUID) ([]models.ActivationX, error) {
	var activations []models.ActivationX
	fmt.Println("Executing SQL Query for partnerID:", partnerID)
	err := r.db.Raw(`
		SELECT 
			a.promo_id, 
			p.title AS promo_name, 
			a.buy, 
			a.current, 
			a.overall, 
			a.created_at, 
			a.updated_at
		FROM activations a
		JOIN promos p ON a.promo_id = p.id
		WHERE p.partner_id = ?
	`, partnerID).Scan(&activations).Error
	fmt.Println("Fetched activations:", activations)
	if err != nil {
		fmt.Println("SQL Error:", err)
		return nil, err
	}
	return activations, nil
}
