package dtos

import (
	"arvan-challenge/services/coupon/internal/db"
	"time"

	"github.com/google/uuid"
)

type CouponDto struct {
	Code         string    `json:"code" validate:"required"`
	ExpireDate   time.Time `json:"expire_date" validate:"required"`
	ChargeAmount int       `json:"charge_amount" validate:"required"`
}

func (cd *CouponDto) ConvertToModel() *db.Coupon {
	id, _ := uuid.NewV7()
	return &db.Coupon{
		Id:           id,
		Code:         cd.Code,
		ExpireDate:   cd.ExpireDate,
		ChargeAmount: cd.ChargeAmount,
	}
}
