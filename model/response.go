package model

import "time"

type PayerResponse struct {
	ID            int       `json:"id" example:"1"`
	Name          *string   `json:"name" example:"Jhon Doe"`
	Email         *string   `json:"email" example:"jhondoe@mail.com"`
	BirthDate     *string   `json:"birth_date" example:"24/07/1992"`
	Phone         *string   `json:"phone" example:"+123456789"`
	Document      *string   `json:"document" xample:"23415162"`
	UserReference *string   `json:"user_reference" example:"12345"`
	Address       Address   `json:"address" gorm:"foreignKey:AddressID;references:ID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
	ID                *string   `json:"id" example:"PAY2323243343543"`
	Amount            float64   `json:"amount" example:"125"`
	Currency          *string   `json:"currency" example:"USD"`
	Country           *string   `json:"country" example:"UY"`
	PaymentMethodID   *string   `json:"payment_method_id" example:"CARD"`
	PaymentMethodFlow *string   `json:"payment_method_flow"`
	OrderNumber       *string   `json:"order_number"`
	Card              Card      `json:"card"`
	CreatedAt         time.Time `json:"created_at"`
}

type OrderResponse struct {
	ID        int               `json:"id"`
	Payer     PayerResponse     `json:"payer"`
	Product   Product           `json:"product" `
	TotalFees int               `json:"total_fees"`
	Payments  []PaymentResponse `json:"payments"`
	Finished  bool              `json:"finished"`
	CreatedAt time.Time         `json:"created_at"`
}

type ProductResponse struct {
	ID          int       `json:"id" example:"1"`
	Name        *string   `json:"name" example:"programacion en C" validate:"nonzero,min=6,max=100"`
	Description *string   `json:"description" example:"Curso de Programacion" validate:"nonzero,min=6,max=100"`
	Amount      float64   `json:"amount" example:"5000.00" validate:"nonzero"`
	Currency    *string   `json:"currency" example:"USD" validate:"nonzero,min=3,max=3"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
