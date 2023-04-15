package dlocal

// SecurePayment
//
// Requires Payer.IP, Payer.DeviceID
// type SecurePaymentRequestBody struct {
// 	Amount            float64 `json:"amount"`
// 	Currency          string  `json:"currency"`
// 	Country           string  `json:"country"`
// 	PaymentMethodID   string  `json:"payment_method_id"`
// 	PaymentMethodFlow string  `json:"payment_method_flow"`
// 	Payer             Payer   `json:"payer"`
// 	Card              Card    `json:"card"`
// 	OrderID           string  `json:"order_id"`
// 	NotificationURL   string  `json:"notification_url"`
// }

// Payment
//
// Use Card's Token saved in 'card' table
type PaymentRequestBody struct {
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Country           string     `json:"country"`
	PaymentMethodID   string     `json:"payment_method_id"`
	PaymentMethodFlow string     `json:"payment_method_flow"`
	Payer             Payer      `json:"payer"`
	Card              SecureCard `json:"card"`
	// OrderID           string     `json:"order_id"`
	// NotificationURL string `json:"notification_url"`
}

// Payment Response
type PaymentResponseBody struct {
	ID                string       `json:"id"`
	Amount            float64      `json:"amount"`
	Currency          string       `json:"currency"`
	Country           string       `json:"country"`
	PaymentMethodID   string       `json:"payment_method_id"`
	PaymentMethodType string       `json:"payment_method_type"`
	PaymentMethodFlow string       `json:"payment_method_flow"`
	Card              CardResponse `json:"card"`
	CreatedDate       string       `json:"created_date"`
	ApprovedDate      string       `json:"approved_date"`
	Status            string       `json:"status"`
	StatusCode        string       `json:"status_code"`
	StatusDetail      string       `json:"status_detail"`
	OrderID           string       `json:"order_id"`
	NotificationUrl   string       `json:"notification_url"`
}
