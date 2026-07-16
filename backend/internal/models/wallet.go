package models

import "time"

type Wallet struct {
	ID           string         `db:"id" json:"id"`
	UserID       string         `db:"user_id" json:"user_id"`
	Currency     WalletCurrency `db:"currency" json:"currency"`
	Balance      int64          `db:"balance" json:"balance"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	WalletNumber int64          `db:"wallet_number" json:"wallet_number"`
}

type WalletCurrency string

const (
	CurrencyUSD WalletCurrency = "USD"
	CurrencyEUR WalletCurrency = "EUR"
	CurrencyTRY WalletCurrency = "TRY"
)
