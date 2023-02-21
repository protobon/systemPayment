package dlocal

// SecurePayment
//
// Requires Payer.IP, Payer.DeviceID
type SecurePaymentRequestBody struct {
	Amount            float64 `json:"amount" validate:"nonzero"`
	Currency          string  `json:"currency" validate:"nonzero"`
	Country           string  `json:"country" validate:"nonzero"`
	PaymentMethodID   string  `json:"payment_method_id" validate:"nonzero"`
	PaymentMethodFlow string  `json:"payment_method_flow" validate:"nonzero"`
	Payer             Payer   `json:"payer" validate:"nonzero"`
	Card              Card    `json:"card" validate:"nonzero"`
	OrderID           string  `json:"order_id" validate:"nonzero"`
	NotificationURL   string  `json:"notification_url" validate:"nonzero"`
}

// Payment
//
// Use Card's Token saved in 'card' table
type PaymentRequestBody struct {
	Amount            float64    `json:"amount" validate:"nonzero"`
	Currency          string     `json:"currency" validate:"nonzero"`
	Country           string     `json:"country" validate:"nonzero"`
	PaymentMethodID   string     `json:"payment_method_id" validate:"nonzero"`
	PaymentMethodFlow string     `json:"payment_method_flow" validate:"nonzero"`
	Payer             Payer      `json:"payer" validate:"nonzero"`
	Card              SecureCard `json:"card" validate:"nonzero"`
	OrderID           string     `json:"order_id" validate:"nonzero"`
	NotificationURL   string     `json:"notification_url" validate:"nonzero"`
}

// Payment Response
type PaymentResponseBody struct {
	ID                *string      `json:"id" validate:"nonzero"`
	Amount            float64      `json:"amount" validate:"nonzero"`
	Currency          *string      `json:"currency" validate:"nonzero"`
	Country           *string      `json:"country" validate:"nonzero"`
	PaymentMethodID   *string      `json:"payment_method_id" validate:"nonzero"`
	PaymentMethodType *string      `json:"payment_method_type" validate:"nonzero"`
	PaymentMethodFlow *string      `json:"payment_method_flow" validate:"nonzero"`
	Card              CardResponse `json:"card" validate:"nonzero"`
	CreatedDate       *string      `json:"created_date" validate:"nonzero"`
	ApprovedDate      *string      `json:"approved_date" validate:"nonzero"`
	Status            *string      `json:"status" validate:"nonzero"`
	StatusCode        *string      `json:"status_code" validate:"nonzero"`
	StatusDetail      *string      `json:"status_detail" validate:"nonzero"`
	OrderID           *string      `json:"order_id" validate:"nonzero"`
	NotificationUrl   *string      `json:"notification_url" validate:"nonzero"`
}
