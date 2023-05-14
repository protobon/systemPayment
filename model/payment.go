package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Payment object
type Payment struct {
	ID                int            `json:"id" gorm:"primaryKey" example:"1"`
	Amount            float64        `json:"amount" example:"5000.00" validate:"nonzero"`
	Currency          *string        `json:"currency" example:"USD" validate:"nonzero,min=3,max=3,uppercase"`
	Country           *string        `json:"country" example:"UY" validate:"nonzero,min=2,max=2,uppercase"`
	PaymentMethodID   *string        `json:"payment_method_id" example:"CARD" validate:"nonzero,min=2,max=4"`
	PaymentMethodFlow *string        `json:"payment_method_flow" example:"DIRECT" validate:"nonzero,min=2,max=10"`
	OrderID           int            `json:"order_id" gorm:"column:order_id" example:"1"  validate:"nonzero"`
	OrderNumber       *string        `json:"order_number" example:"657434343"  validate:"nonzero"`
	CardID            int            `json:"card_id" gorm:"column:card_id" example:"1"  validate:"nonzero"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `json:"-"`
}

func (Payment) TableName() string {
	return "payment"
}

// Save payment from dlocal's payment response
func (p *Payment) SavePaymentFromResponse(db *gorm.DB, response map[string]interface{}) (int, error) {
	amount, _ := response["amount"].(float64)
	currency, _ := response["currency"].(string)
	country, _ := response["country"].(string)
	payment_method_id, _ := response["payment_method_id"].(string)
	payment_method_flow, _ := response["payment_method_flow"].(string)
	order_number, _ := response["order_id"].(string)

	p.Amount = amount
	p.Currency = &currency
	p.Country = &country
	p.PaymentMethodID = &payment_method_id
	p.PaymentMethodFlow = &payment_method_flow
	p.OrderNumber = &order_number
	p.CreatedAt = time.Now()

	return p.QCreatePayment(db)
}

// QCreatePayment - Insert into Payment
//
// Inserts new Payment
func (p *Payment) QCreatePayment(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(p); err != nil {
		return 400, err
	}

	p.CreatedAt = time.Now()
	// Create Payment
	if err = db.Create(p).Error; err != nil {
		return 500, err
	}

	return 200, nil
}

// QGetPayments - Get payments from order
func (p *Payment) QGetPayments(db *gorm.DB) ([]Payment, int, error) {
	var payments []Payment
	if err := db.Table("payment").Select("*").Where("order_id=?", p.OrderID).Scan(&payments).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return payments, 404, err
		default:
			return payments, 500, err
		}
	}

	return payments, 200, nil
}

// QGetPayment - Get payment from id
func (p *Payment) QGetPayment(db *gorm.DB) (int, error) {
	if err := db.Where("id = ?", p.ID).First(&p).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}

// Get all payments (optional order_id)
func (p *Payment) QGetAllPayments(db *gorm.DB, start int, count int, order_id int) ([]Payment, int, error) {
	var payments []Payment
	if order_id != 0 {
		if err := db.Table("payment").Where("order_id=?", order_id).Select("*").
			Order("created_at desc").Limit(count).Offset(start).Scan(&payments).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return payments, 404, err
			default:
				return payments, 500, err
			}
		}
	} else {
		if err := db.Table("payment").Select("*").Order("created_at desc").
			Limit(count).Offset(start).Scan(&payments).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return payments, 404, err
			default:
				return payments, 500, err
			}
		}
	}

	return payments, 200, nil
}
