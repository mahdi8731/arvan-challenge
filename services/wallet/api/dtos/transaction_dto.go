package dtos

import (
	"arvan-challenge/services/wallet/internal/db"
	"time"

	"github.com/google/uuid"
)

type TransactionDto struct {
	Description string `json:"description" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Amount      int    `json:"amount" validate:"required"`
}

func (t *TransactionDto) ConvertToModel() (*db.Transaction, *db.Wallet) {
	tid, _ := uuid.NewV7()
	wid, _ := uuid.NewV7()
	tt := time.Now().UTC()
	return &db.Transaction{
			Id:          tid,
			Date:        tt,
			Description: t.Description,
			Amount:      t.Amount,
		}, &db.Wallet{
			Id:           wid,
			PhoneNumber:  t.PhoneNumber,
			LastModefied: tt,
		}
}
