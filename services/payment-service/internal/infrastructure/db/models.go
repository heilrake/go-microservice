package infrastructure

import "time"

// PaymentIntent represents the payment_intents table in the database
type PaymentIntent struct {
	ID                    string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	TripID                string    `gorm:"column:trip_id;type:varchar(255);not null;uniqueIndex"`
	UserID                string    `gorm:"column:user_id;type:varchar(255);not null;index"`
	StripePaymentIntentID string    `gorm:"column:stripe_payment_intent_id;type:varchar(255);not null"`
	Amount                int64     `gorm:"column:amount;not null"`
	Currency              string    `gorm:"column:currency;type:varchar(10);not null;default:'usd'"`
	Status                string    `gorm:"column:status;type:varchar(50);not null;default:'authorized'"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (PaymentIntent) TableName() string {
	return "payment_intents"
}
