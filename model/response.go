package model

import "time"

type PayerResponse struct {
	ID            int             `json:"id" example:"1"`
	Name          *string         `json:"name" example:"Jhon Doe"`
	Email         *string         `json:"email" example:"jhondoe@mail.com"`
	BirthDate     *string         `json:"birth_date" example:"24/07/1992"`
	Phone         *string         `json:"phone" example:"+123456789"`
	Document      *string         `json:"document" xample:"23415162"`
	UserReference *string         `json:"user_reference" example:"12345"`
	Address       AddressResponse `json:"address"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type AddressResponse struct {
	ID        int       `json:"id" example:"1"`
	State     *string   `json:"state" example:"Rio de Janeiro"`
	City      *string   `json:"city" example:"Volta Redonda"`
	ZipCode   *string   `json:"zip_code" example:"27275-595"`
	Street    *string   `json:"street" example:"Servidão B-1"`
	Number    *string   `json:"number" example:"1106"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentResponse struct {
	ID                *string // Dlocal payment.id
	Amount            float64
	Currency          *string
	Country           *string
	PaymentMethodID   *string
	PaymentMethodFlow *string
	OrderNumber       *string // Dlocal order_id
	Card              Card
	CreatedAt         time.Time
}

type OrderResponse struct {
	ID        int               `json:"id"`
	Payer     PayerResponse     `json:"payer"`
	Product   Product           `json:"product" `
	TotalFees int               `json:"total_fees"`
	Payments  []PaymentResponse `json:"payments"`
	Finished  bool              `json:"finished"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
