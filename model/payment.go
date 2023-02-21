package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Payment object
type Payment struct {
	ID                int            `gorm:"primaryKey" example:"1"`
	Amount            float64        `example:"5000.00" validate:"nonzero"`
	Currency          *string        `example:"USD" validate:"nonzero,min=3,max=3"`
	Country           *string        `example:"UY" validate:"nonzero,min=2,max=2"`
	PaymentMethodID   *string        `example:"CARD" validate:"nonzero,min=2,max=4"`
	PaymentMethodFlow *string        `example:"DIRECT" validate:"nonzero,min=2,max=10"`
	OrderID           int            `gorm:"column:order_id" example:"1"  validate:"nonzero"`
	OrderNumber       *string        `example:"657434343"  validate:"nonzero"`
	CardID            int            `gorm:"column:card_id" example:"1"  validate:"nonzero"`
	CreatedAt         time.Time      `json:"-" swaggerignore:"true"`
	UpdatedAt         time.Time      `json:"-" swaggerignore:"true"`
	DeletedAt         gorm.DeletedAt `json:"-" swaggerignore:"true"`
}

func (Payment) TableName() string {
	return "payment"
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

func (p *Payment) QGetPayments(db *gorm.DB, start int, count int) ([]Payment, int, error) {
	var payments []Payment
	if err := db.Table("payment").Select("*").Scan(&payments).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return payments, 404, err
		default:
			return payments, 500, err
		}
	}

	return payments, 200, nil
}

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

func (p *Payment) QUpdatePayment(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(p); err != nil {
		return 400, err
	}

	p.UpdatedAt = time.Now()
	if err = db.Model(&p).Updates(p).Error; err != nil {
		return 500, err
	}
	return 200, nil
}
