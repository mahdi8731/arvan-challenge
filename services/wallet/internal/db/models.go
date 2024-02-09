package db

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id          uuid.UUID `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"Description"`
	Amount      int       `json:"amount"`
}

type Wallet struct {
	Id           uuid.UUID `json:"id"`
	LastModefied time.Time `json:"last_modified"`
	PhoneNumber  string    `json:"phone_number" `
	Inventory    int       `json:"inventory" `
}
