package model

type PayerRequest struct {
	Name          *string        `json:"name" example:"Jhon Doe"`
	Email         *string        `json:"email" example:"jhondoe@mail.com"`
	BirthDate     *string        `json:"birth_date" example:"24/07/1992"`
	Phone         *string        `json:"phone" example:"+123456789"`
	Document      *string        `json:"document" example:"23415162"`
	UserReference *string        `json:"user_reference" example:"12345"`
	Address       AddressRequest `json:"address" gorm:"foreignKey:AddressID;references:ID"`
}

type AddressRequest struct {
	ID      int     `json:"id" example:"1"`
	State   *string `json:"state" example:"Rio de Janeiro"`
	City    *string `json:"city" example:"Volta Redonda"`
	ZipCode *string `json:"zip_code" example:"27275-595"`
	Street  *string `json:"street" example:"Servidão B-1"`
	Number  *string `json:"number" example:"1106"`
}

type CardRequest struct {
	Token *string `json:"token" validate:"nonzero"`
	Last4 *string `json:"last_4" example:"1234" validate:"nonzero"`
	Brand *string `json:"brand" example:"Visa" validate:"nonzero"`
}

type PaymentRequest struct {
	Amount            float64 `json:"amount" example:"125"`
	Currency          *string `json:"currency" example:"USD"`
	Country           *string `json:"country" example:"UY"`
	PaymentMethodID   *string `json:"payment_method_id" example:"CARD"`
	PaymentMethodFlow *string `json:"payment_method_flow" example:"DIRECT"`
	OrderNumber       *string `json:"order_number" example:"657434343"`
}

type OrderRequest struct {
	ProductID int     `json:"product_id" example:"1"`
	Currency  *string `json:"currency" validate:"nonzero,min=3,max=3"`
	TotalFees int     `json:"total_fees" validate:"nonzero" example:"3"`
}

type ProductRequest struct {
	Name        *string `json:"name" example:"programacion en C" validate:"nonzero,min=6,max=100"`
	Description *string `json:"description" example:"Curso de Programacion" validate:"nonzero,min=6,max=100"`
	Amount      float64 `json:"amount" example:"5000.00" validate:"nonzero"`
	Currency    *string `json:"currency" example:"USD" validate:"nonzero,min=3,max=3"`
}
