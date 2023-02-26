package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Product example
type Product struct {
	ID          int            `json:"id" gorm:"primaryKey" example:"1"`
	Name        *string        `json:"name" example:"programacion en C" validate:"nonzero,min=6,max=100"`
	Description *string        `json:"description" example:"Curso de Programacion" validate:"nonzero,min=6,max=100"`
	Amount      float64        `json:"amount" example:"5000.00" validate:"nonzero"`
	Currency    *string        `json:"currency" example:"USD" validate:"nonzero,min=3,max=3"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func (Product) TableName() string {
	return "product"
}

func ProductExists(db *gorm.DB, id int) (bool, error) {
	var p Product
	if err := db.Table("product").Select("id").Where("id=?", id).First(&p).Error; err != nil {
		return false, err
	}
	return true, nil
}

// QCreateproduct - Insert into product
//
// Inserts new product
func (p *Product) QCreateProduct(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(p); err != nil {
		return 400, err
	}

	p.CreatedAt = time.Now()
	// Create product
	if err = db.Create(p).Error; err != nil {
		return 500, err
	}

	return 200, nil
}

func (p *Product) QGetProducts(db *gorm.DB, start int, count int) ([]Product, int, error) {
	var products []Product
	if err := db.Table("product").Select("*").Scan(&products).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return products, 404, err
		default:
			return products, 500, err
		}
	}

	return products, 200, nil
}

func (p *Product) QGetProduct(db *gorm.DB) (int, error) {
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

func (p *Product) QUpdateProduct(db *gorm.DB) (int, error) {
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
