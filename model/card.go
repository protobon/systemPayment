package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Card example
type Card struct {
	ID        int            `json:"id" gorm:"primaryKey" example:"1"`
	PayerID   int            `json:"payer_id" gorm:"column:payer_id" example:"1"  validate:"nonzero,min=1"`
	Token     *string        `json:"token" validate:"nonzero"`
	Last4     *string        `json:"last_4" gorm:"column:last_4" example:"1234" validate:"nonzero,min=4,max=4"`
	Brand     *string        `json:"brand" example:"Visa" validate:"nonzero"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Card) TableName() string {
	return "card"
}

// QCreateCard
//
// Inserts new Card
func (c *Card) QCreateCard(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(c); err != nil {
		return 400, err
	}
	c.CreatedAt = time.Now()

	if err = db.Create(c).Error; err != nil {
		return 500, err
	}
	return 200, nil
}

// QGetCards
//
// Get Payer's Secured Cards (match Card.PayerID)
func (c *Card) QGetCards(db *gorm.DB, payer_id int) ([]Card, int, error) {
	var cards []Card
	if err := db.Table("card").Where("payer_id=?", c.PayerID).Select("*").Scan(&cards).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return cards, 404, err
		default:
			return cards, 500, err
		}
	}

	return cards, 200, nil
}

// QGetCard
//
// Get one Card from Card.ID and Card.PayerID
func (c *Card) QGetCard(db *gorm.DB) (int, error) {
	if err := db.Where("id = ?", c.ID).First(&c).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}
