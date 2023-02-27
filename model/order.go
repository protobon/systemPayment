package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

type OrderMetadata struct {
	Order   Order   `json:"order"`
	Product Product `json:"product"`
}

// Order object
type Order struct {
	ID        int            `json:"id" gorm:"primaryKey" example:"1"`
	PayerID   int            `json:"payer_id" gorm:"column:payer_id" example:"1"  validate:"nonzero"`
	ProductID int            `json:"product_id" example:"1"  validate:"nonzero"`
	Product   Product        `json:"product"`
	TotalFees int            `json:"total_fees" example:"3"  validate:"nonzero,min=1,max=24"`
	Payments  []Payment      `json:"payments"`
	Finished  bool           `json:"finished" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Order) TableName() string {
	return "order"
}

func OrderExists(db *gorm.DB, id int) (bool, error) {
	var o Order
	if err := db.Table("order").Select("id").Where("id=?", id).First(&o).Error; err != nil {
		return false, err
	}
	return true, nil
}

func PreloadOrder(db *gorm.DB, order_id int, payer_id int) (*Order, error) {
	var o *Order
	if err := db.Table("order").Select("id, payer_id").Where("id=? AND payer_id=?", order_id, payer_id).
		First(&o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

// QCreateOrder - Insert into Order
//
// Inserts new Order
func (o *Order) QCreateOrder(db *gorm.DB) (int, error) {
	var err error
	if t, err := PayerExists(db, o.PayerID); !t {
		return 404, err
	}

	if t, err := ProductExists(db, o.ProductID); !t {
		return 404, err
	}

	var o_req = OrderRequest{TotalFees: o.TotalFees}
	if err = validator.Validate(o_req); err != nil {
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
	if err := db.Model(&Order{}).Preload("Product").Find(&orders).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return orders, 404, err
		default:
			return orders, 500, err
		}
	}
	for idx, order := range orders {
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
	if err := db.Table("order").Preload("Product").Where("id=?", o.ID).First(&o).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
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
