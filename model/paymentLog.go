package model

type PaymentMetadata struct {
	PaymentID int
	PayerID   int
	Success   bool
	Body      Payment
}

// PaymentLog example
type PaymentLog struct {
	PayerID  int
	Payments []PaymentMetadata
}

func (PaymentMetadata) TableName() string {
	return "payment_metadata"
}
