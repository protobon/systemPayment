package model

import (
	"database/sql"
	"time"
)

type SecureCardObject struct {
	ID        int       `json:"id"`
	PayerID   int       `json:"payer_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type CardObject struct {
	ID              int       `json:"-"`
	HolderName      string    `json:"holder_name"`
	Number          string    `json:"number"`
	CVV             string    `json:"cvv"`
	ExpirationMonth uint8     `json:"expiration_month"`
	ExpirationYear  uint16    `json:"expiration_year"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DlocalSecureCard struct {
	ID        int       `json:"-"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type DlocalCard struct {
	ID              int       `json:"-"`
	HolderName      string    `json:"holder_name"`
	Number          string    `json:"number"`
	CVV             string    `json:"cvv"`
	ExpirationMonth uint8     `json:"expiration_month"`
	ExpirationYear  uint16    `json:"expiration_year"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

// QCreateSecureCard - Inserta en tabla 'secure_card'
func (card *SecureCardObject) QCreateSecureCard(db *sql.DB) error {
	card.CreatedAt = time.Now()
	err := db.QueryRow(
		"INSERT INTO secure_card(payer_id, token, created_at) VALUES($1, $2, $3) RETURNING id",
		card.PayerID, card.Token, card.CreatedAt).Scan(&card.ID)

	if err != nil {
		return err
	}

	return nil
}
