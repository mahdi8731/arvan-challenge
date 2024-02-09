package db

import (
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	Id           uuid.UUID `json:"id"`
	ExpireDate   time.Time `json:"expire_date"`
	Code         string    `json:"code"`
	ChargeAmount int       `json:"charge_amount"`
	AllowedTimes int       `json:"allowed_times"`
}
