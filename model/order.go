package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Order example
type Order struct {
	ID         int            `gorm:"primaryKey" example:"1" validate:"nonzero"`
	PayerID    int            `gorm:"column:payer_id" example:"1"  validate:"nonzero"`
	ProductID  int            `gorm:"column:product_id" example:"1"  validate:"nonzero"`
	TotalFees  uint           `example:"3"  validate:"nonzero"`
	CreatedAt  time.Time      `json:"-" swaggerignore:"true"`
	UpdatedAt  time.Time      `json:"-" swaggerignore:"true"`
	DeletedAt  gorm.DeletedAt `json:"-" swaggerignore:"true"`
	FinishedAt time.Time      `json:"-" swaggerignore:"true"`
	Finished   bool           `gorm:"default:false"`
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
	if err := db.Table("order").Select("*").Scan(&orders).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return orders, 404, err
		default:
			return orders, 500, err
		}
	}

	return orders, 200, nil
}

func (o *Order) QGetOrder(db *gorm.DB) (int, error) {
	if err := db.Where("id = ?", o.ID).First(&o).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
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
