package model

import (
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
		log.Error("QCreateOrder - ", err)
		return 400, err
	}

	var o_req = OrderRequest{
		ProductID: o.ProductID,
		TotalFees: o.TotalFees,
		Currency:  o.Currency,
	}
	if err = validator.Validate(o_req); err != nil {
		log.Error("QCreateOrder - ", err)
		return 400, err
	}

	var product = Product{ID: o.ProductID}
	code, err := product.QGetProduct(db)
	if err != nil {
		return code, err
	}

	o.OrderId = uuid.New().String()
	o.CreatedAt = time.Now()
	o.NextPayment = time.Now()
	o.CurrentFee = 1

	// Create Order
	if err = db.Create(o).Error; err != nil {
		log.Error("QCreateOrder - ", err)
		return 400, err
	}

	return 200, nil
}

func (o *Order) QGetOrders(db *gorm.DB, start int, count int, payer_id int) ([]Order, int, error) {
	var orders []Order
	if payer_id != 0 {
		if err := db.Model(&Order{}).Where("payer_id=?", payer_id).Preload("Product").Limit(count).
			Offset(start).Find(&orders).Error; err != nil {
			log.Error("QGetOrders - ", err)
			return orders, 400, err
		}
	} else {
		if err := db.Model(&Order{}).Preload("Product").Limit(count).Offset(start).
			Find(&orders).Error; err != nil {
			log.Error("QGetOrders - ", err)
			return orders, 400, err
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
		log.Error("QGetOrder - ", err)
		return 400, err
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
		log.Error("QUpdateOrder - ", err)
		return 400, err
	}

	o.UpdatedAt = time.Now()
	if err = db.Model(&o).Updates(o).Error; err != nil {
		log.Error("QUpdateOrder - ", err)
		return 400, err
	}
	return 200, nil
}

// Fetches one order by ID only with the necessary data for making a payment
func (o *Order) GetOrderForPayment(db *gorm.DB) (int, error) {
	if err := db.Table("order").Where("id=?", o.ID).Where("finished=?", false).First(&o).Error; err != nil {
		log.Error("GetOrderForPayment - ", err)
		return 400, err
	}
	return 200, nil
}

// Handles order after successful payment
func (o *Order) PaymentSuccessful(db *gorm.DB) (int, error) {
	if o.CurrentFee == o.TotalFees {
		o.Finished = true
		o.Auto = false
	} else {
		o.CurrentFee++
		_ = o.NextPayment.AddDate(0, 1, 0)
	}
	return o.QUpdateOrder(db)
}
