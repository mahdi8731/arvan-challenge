package db

import (
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	Id           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	ExpireDate   time.Time `json:"expire_date"`
	ChargeAmount int       `json:"charge_amount"`
}
