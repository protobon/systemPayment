package model

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Order object
type Order struct {
	ID          int            `json:"id" gorm:"primaryKey" example:"1"`
	Amount      float64        `json:"amount"`
	OrderId     string         `json:"order_id"`
	Currency    *string        `json:"currency" example:"USD" validate:"nonzero"`
	PayerID     int            `json:"payer_id" gorm:"column:payer_id" example:"1"  validate:"nonzero"`
	ProductID   int            `json:"product_id" example:"1"  validate:"nonzero"`
	Product     Product        `json:"product"`
	TotalFees   int            `json:"total_fees" example:"3"  validate:"nonzero,min=1,max=24"`
	CurrentFee  int            `json:"current_fee" example:"1"`
	Auto        bool           `json:"-"`
	NextPayment time.Time      `json:"next_payment"`
	Payments    []Payment      `json:"payments"`
	Finished    bool           `json:"finished" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (Order) TableName() string {
	return "order"
}

// QCreateOrder - Insert into Order
//
// Inserts new Order
func (o *Order) QCreateOrder(db *gorm.DB) (int, error) {
	var err error
	if t, err := PayerExists(db, o.PayerID); !t {
		return 404, err
	}

	var o_req = OrderRequest{
		ProductID: o.ProductID,
		TotalFees: o.TotalFees,
		Currency:  o.Currency,
	}
	if err = validator.Validate(o_req); err != nil {
		return 400, err
	}

	var product = Product{ID: o.ProductID}
	code, err := product.QGetProduct(db)
	if err != nil {
		return code, err
	}

	o.OrderId = uuid.New().String()
	o.CreatedAt = time.Now()
	o.CurrentFee = 1
	// Create Order
	if err = db.Create(o).Error; err != nil {
		return 500, err
	}

	return 200, nil
}

func (o *Order) QGetOrders(db *gorm.DB, start int, count int, payer_id int) ([]Order, int, error) {
	var orders []Order
	if payer_id != 0 {
		if err := db.Model(&Order{}).Where("payer_id=?", payer_id).Preload("Product").Limit(count).
			Offset(start).Find(&orders).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return orders, 404, err
			default:
				return orders, 500, err
			}
		}
	} else {
		if err := db.Model(&Order{}).Preload("Product").Limit(count).Offset(start).
			Find(&orders).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return orders, 404, err
			default:
				return orders, 500, err
			}
		}
	}

	for idx, order := range orders {
		p := Payment{OrderID: order.ID}
		payments, code, err := p.QGetPayments(db)
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
	payments, code, err := p.QGetPayments(db)
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
