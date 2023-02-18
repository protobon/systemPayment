package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Dummy example
type Dummy struct {
	ID        int            `gorm:"primaryKey" example:"1"`
	Name      *string        `example:"desktop chair" validate:"nonzero"`
	Price     float64        `example:"299.99" validate:"nonzero"`
	CreatedAt time.Time      `swaggerignore:"true"`
	UpdatedAt time.Time      `swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"-" swaggerignore:"true"`
}

func (Dummy) TableName() string {
	return "dummy"
}

func (d *Dummy) QCreateDummy(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(d); err != nil {
		return 400, err
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	if err = db.Create(d).Error; err != nil {
		return 500, err
	}
	return 200, nil
}

func (d *Dummy) QGetDummies(db *gorm.DB, start int, count int) ([]Dummy, int, error) {
	var dummies []Dummy
	if err := db.Table("dummy").Select("*").Scan(&dummies).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return dummies, 404, err
		default:
			return dummies, 500, err
		}
	}

	return dummies, 200, nil
}

func (d *Dummy) QGetDummy(db *gorm.DB) (int, error) {
	if err := db.Where("id = ?", d.ID).First(&d).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}

func (d *Dummy) QUpdateDummy(db *gorm.DB) (int, error) {
	var err error
	d.UpdatedAt = time.Now()
	if err = db.Model(&d).Updates(d).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}
