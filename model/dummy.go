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

func (d *Dummy) QCreateDummy(db *gorm.DB) error {
	var err error
	if err = validator.Validate(d); err != nil {
		return err
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	err = db.Create(d).Error
	return err
}

func (d *Dummy) QGetDummies(db *gorm.DB, start int, count int) ([]Dummy, error) {
	var dummies []Dummy
	if err := db.Table("dummy").Select("*").Scan(&dummies).Error; err != nil {
		return nil, err
	}

	return dummies, nil
}

func (d *Dummy) QGetDummy(db *gorm.DB) error {
	err := db.Where("id = ?", d.ID).First(&d).Error
	return err
}

func (d *Dummy) QUpdateDummy(db *gorm.DB) error {
	var err error
	d.UpdatedAt = time.Now()
	err = db.Model(&d).Updates(d).Error
	return err
}
