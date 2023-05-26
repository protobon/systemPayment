package model

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Card's token
type Token struct {
	Token string `json:"token" default:""`
}

// Card example
type Card struct {
	ID        int            `json:"id" gorm:"primaryKey" example:"1"`
	PayerID   int            `json:"payer_id" gorm:"column:payer_id" example:"1"  validate:"nonzero,min=1"`
	CardId    *string        `json:"card_id" validate:"nonzero"`
	Last4     *string        `json:"last_4" gorm:"column:last_4" example:"1234" validate:"nonzero,min=4,max=4"`
	Brand     *string        `json:"brand" example:"Visa" validate:"nonzero"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Card) TableName() string {
	return "card"
}

// Save card from dlocal's payment response
func (c *Card) SaveCardFromResponse(db *gorm.DB, response map[string]interface{}) (int, error) {
	card, _ := response["card"].(map[string]interface{})
	cardId, _ := card["card_id"].(string)
	last4, _ := card["last4"].(string)
	brand, _ := card["brand"].(string)

	c.CardId = &cardId
	c.Last4 = &last4
	c.Brand = &brand
	c.CreatedAt = time.Now()

	return c.QCreateCard(db)
}

// QCreateCard
//
// Inserts new Card
func (c *Card) QCreateCard(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(c); err != nil {
		log.Error("QCreateCard - ", err)
		return 400, err
	}
	c.CreatedAt = time.Now()

	if err = db.Create(c).Error; err != nil {
		log.Error("QCreateCard - ", err)
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
		log.Error("QGetCards - ", err)
		switch err {
		case gorm.ErrRecordNotFound:
			return cards, 200, err
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
		log.Error("Get Card - " + err.Error())
		return 400, err
	}
	return 200, nil
}
