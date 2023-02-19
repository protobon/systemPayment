package dlocal

// Payment
//
// Requires Payer.IP, Payer.DeviceID
type Payment struct {
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Country           string  `json:"country"`
	PaymentMethodID   string  `json:"payment_method_id"`
	PaymentMethodFlow string  `json:"payment_method_flow"`
	Payer             Payer   `json:"payer"`
	Card              Card    `json:"card"`
	OrderID           string  `json:"order_id"`
	NotificationURL   string  `json:"notification_url"`
}

// SecurePayment
//
// Use Card's Token saved in 'card' table
type SecurePayment struct {
	Amount            float64    `json:"amount"`
	Currency          string     `json:"currency"`
	Country           string     `json:"country"`
	PaymentMethodID   string     `json:"payment_method_id"`
	PaymentMethodFlow string     `json:"payment_method_flow"`
	Payer             Payer      `json:"payer"`
	Card              SecureCard `json:"card"`
	OrderID           string     `json:"order_id"`
	NotificationURL   string     `json:"notification_url"`
}
