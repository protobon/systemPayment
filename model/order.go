package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Order object
type Order struct {
	ID        int            `json:"id" gorm:"primaryKey" example:"1"`
	PayerID   int            `json:"-" gorm:"column:payer_id" example:"1"  validate:"nonzero"`
	Payer     Payer          `json:"payer"`
	ProductID int            `json:"-" gorm:"column:product_id" example:"1"  validate:"nonzero"`
	Product   Product        `json:"product"`
	TotalFees int            `json:"total_fees" example:"3"  validate:"nonzero,min=1,max=24"`
	Payments  []Payment      `json:"payments"`
	Finished  bool           `json:"-" gorm:"default:false" swaggerignore:"true"`
	CreatedAt time.Time      `json:"-" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"-" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
}

func (Order) TableName() string {
	return "order"
}

// QCreateOrder - Insert into Order
//
// Inserts new Order
func (o *Order) QCreateOrder(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(o); err != nil {
		return 400, err
	}

	o.CreatedAt = time.Now()
	// Create Order
	if err = db.Create(o).Error; err != nil {
		return 500, err
	}

	return 200, nil
}

func (o *Order) QGetOrders(db *gorm.DB, start int, count int) ([]Order, int, error) {
	var orders []Order
	if err := db.Preload("Product").Where("payer_id=?", o.PayerID).Find(&orders).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return orders, 404, err
		default:
			return orders, 500, err
		}
	}
	for idx, order := range orders {
		order.Payer.ID = order.PayerID
		if code, err := order.Payer.QGetPayer(db); err != nil {
			return orders, code, err
		}
		p := Payment{OrderID: order.ID}
		payments, code, err := p.QGetPayments(db, 0, 10)
		if err != nil {
			return orders, code, err
		}
		order.Payments = payments
		orders[idx] = order
	}
	return orders, 200, nil
}

func (o *Order) QGetOrder(db *gorm.DB) (int, error) {
	if err := db.Preload("Product").Where("id=?", o.ID).First(&o).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	o.Payer.ID = o.PayerID
	if code, err := o.Payer.QGetPayer(db); err != nil {
		return code, err
	}
	p := Payment{OrderID: o.ID}
	payments, code, err := p.QGetPayments(db, 0, 10)
	if err != nil {
		return code, err
	}
	o.Payments = payments
	return 200, nil
}

func (o *Order) QUpdateOrder(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(o); err != nil {
		return 400, err
	}

	o.UpdatedAt = time.Now()
	if err = db.Model(&o).Updates(o).Error; err != nil {
		return 500, err
	}
	return 200, nil
}
