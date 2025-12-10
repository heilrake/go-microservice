package domain

import "time"

type RideFareModel struct {
	ID                string
	UserID            string
	PackageSlug       string // e.g., "standard", "premium"
	TotalPriceInCents float64
	ExpiresAt         time.Time
}
