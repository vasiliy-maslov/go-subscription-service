package model

import (
	"github.com/google/uuid"
	"time"
)

// Subscription представляет одну запись о подписке
type Subscription struct {
	ID          uuid.UUID  `db:"id"           json:"id"`
	UserID      uuid.UUID  `db:"user_id"      json:"user_id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	Price       int        `db:"price"        json:"price"`
	StartDate   time.Time  `db:"start_date"   json:"start_date"`
	EndDate     *time.Time `db:"end_date"     json:"end_date,omitempty"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"   json:"updated_at"`
}
