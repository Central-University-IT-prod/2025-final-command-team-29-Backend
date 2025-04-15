package models

import (
	"time"

	"github.com/google/uuid"
)

// Stat представляет статистику бизнеса по промо-акциям.
//
// @Description Структура, содержащая основные бизнес-метрики использования промо-кодов.
type Stat struct {
	// Количество активных промо-акций, которые в данный момент доступны клиентам.
	// Промо-акция считается активной, если у нее `is_active = true`.
	ActivPromos int `json:"active_promos_count" example:"5"`

	// Общее количество покупок, совершенных с использованием промо-кодов.
	// Считается суммарное количество активаций промо, приведших к покупке.
	BuyPromos int `json:"buy_promos_count" example:"120"`

	// Количество уникальных пользователей, которые использовали промо-коды.
	// Каждый клиент считается только один раз, независимо от количества активаций.
	ActiveUsers int `json:"active_users_count" example:"85"`

	// Количество промо-акций, достигших лимита использования.
	// Например, если промо-код имеет ограничение `usage_limit = 100`, то он больше не учитывается после достижения лимита.
	LimitPromos int `json:"limit_promos_count" example:"10"`

	// Количество покупок, совершенных с промо-кодами за последние 24 часа.
	// Полезно для оценки текущей активности пользователей.
	DailyBuyPromo int `json:"daily_buy_promo_count" example:"30"`

	// Количество уникальных пользователей, активировавших промо-коды за последние 24 часа.
	// Показывает ежедневную вовлеченность клиентов в акции.
	DailyActiveUsers int `json:"daily_active_users_count" example:"25"`

	// Количество клиентов, которые использовали промо-коды повторно в последние 30 дней.
	// Этот показатель помогает измерить уровень удержания клиентов.
	Retention30d int `json:"retention_30d" example:"40"`

	// Общее количество зарегистрированных клиентов у партнера.
	// Полезно для оценки охвата бизнеса и анализа конверсии в использование промо-кодов.
	TotalClients int `json:"total_users" example:"500"`

	// Конверсия пользователей в использование промо-кодов (в процентах).
	// Рассчитывается как `(ActiveUsers / TotalClients) * 100`.
	ConversionRate float64 `json:"conversion_rate" example:"17.0"`

	// Количество топ-клиентов, которые использовали промо-коды чаще всего.
	// Этот показатель можно использовать для определения лояльных клиентов.
	TopClients int `json:"top_users" example:"10"`
}
type StatPromo struct {
	BuyPromos     int `json:"buy_promos_count"`
	Users         int `json:"users_count"`
	DailyBuyPromo int `json:"daily_buy_promo_count"`
	DailyUsers    int `json:"daily_users_count"`
}

type ActivationX struct {
	PromoName string    `gorm:"column:promo_name"`
	PromoID   uuid.UUID `gorm:"column:promo_id"`
	Buy       int       `gorm:"column:buy"`
	Current   int       `gorm:"column:current"`
	Overall   int       `gorm:"column:overall"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
