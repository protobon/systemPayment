package model

import "time"

type PaymentObject struct {
	ID                int         `json:"id"`
	Amount            float64     `json:"amount"`
	Currency          string      `json:"currency"`
	Country           string      `json:"country"`
	PaymentMethodID   string      `json:"payment_method_id"`
	PaymentMethodFlow string      `json:"payment_method_flow"`
	Payer             PayerObject `json:"payer"`
	Card              CardObject  `json:"card"`
	OrderID           string      `json:"order_id"`
	NotificationURL   string      `json:"notification_url"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

type DlocalPayment struct {
	ID                int         `json:"-"`
	Amount            float64     `json:"amount"`
	Currency          string      `json:"currency"`
	Country           string      `json:"country"`
	PaymentMethodID   string      `json:"payment_method_id"`
	PaymentMethodFlow string      `json:"payment_method_flow"`
	Payer             PayerObject `json:"payer"`
	Card              CardObject  `json:"card"`
	OrderID           string      `json:"order_id"`
	NotificationURL   string      `json:"notification_url"`
	CreatedAt         time.Time   `json:"-"`
	UpdatedAt         time.Time   `json:"-"`
}
