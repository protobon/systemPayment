package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Payer example
type Payer struct {
	ID            int            `json:"id" gorm:"primaryKey" example:"1"`
	Name          *string        `json:"name" example:"Jhon Doe" validate:"nonzero,min=3,max=100"`
	Email         *string        `json:"email" example:"jhondoe@mail.com" validate:"nonzero,min=8,max=100"`
	BirthDate     *string        `json:"birth_date" example:"24/07/1992" validate:"nonzero"`
	Phone         *string        `json:"phone" example:"+123456789" validate:"nonzero"`
	Document      *string        `json:"document" example:"23415162" validate:"nonzero"`
	UserReference string         `json:"user_reference"`
	Address       Address        `json:"address" gorm:"foreignKey:PayerID;references:ID" validate:"nonzero"`
	AddressID     int            `json:"-"`
	Country       *string        `json:"country" example:"UY" validate:"nonzero,min=2,max=2"`
	CardID        int            `json:"card_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

type Address struct {
	ID        int            `json:"-" gorm:"primaryKey" example:"1"`
	PayerID   int            `json:"-" gorm:"column:payer_id" example:"1"`
	State     *string        `json:"state" example:"Rio de Janeiro" validate:"nonzero"`
	City      *string        `json:"city" example:"Volta Redonda" validate:"nonzero"`
	ZipCode   *string        `json:"zip_code" example:"27275-595" validate:"nonzero"`
	Street    *string        `json:"street" example:"Servidão B-1" validate:"nonzero"`
	Number    *string        `json:"number" example:"1106" validate:"nonzero"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Payer) TableName() string {
	return "payer"
}

func (Address) TableName() string {
	return "address"
}

func PayerExists(db *gorm.DB, id int) (bool, error) {
	var p Payer
	if err := db.Table("payer").Select("id").Where("id=?", id).First(&p).Error; err != nil {
		return false, err
	}
	return true, nil
}

func PreloadPayer(db *gorm.DB, id int) (*Payer, error) {
	var p *Payer
	if err := db.Table("payer").Select("id, card_id").Where("id=?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

// QGetPayerFromEmail - Get payer from email
//
// Filter by Payer.Email, returns Payer{} or ErrRecordNotFound
func (p *Payer) QGetPayerFromEmail(db *gorm.DB) error {
	err := db.Where("email = ?", p.Email).First(&p).Error
	return err
}

// QCreatePayer - Insert into payer
//
// Inserts new Payer + Address
func (p *Payer) QCreatePayer(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(p); err != nil {
		return 400, err
	}
	if err = validator.Validate(p.Address); err != nil {
		return 400, err
	}

	p.CreatedAt = time.Now()
	// Create Payer
	if err = db.Raw(`INSERT INTO payer(name, email, country, birth_date, phone, document, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`, p.Name, p.Email, p.Country, p.BirthDate, p.Phone,
		p.Document, p.CreatedAt).Scan(&p.ID).Error; err != nil {
		return 500, err
	}
	str_payer_id := strconv.Itoa(p.ID)
	p.UserReference = fmt.Sprintf("%05s", str_payer_id)

	// Insert Payer's address
	return p.Address.QCreateAddress(db, p)
}

// QCreateAddress - Insert into address
//
// Inserts new Address linked to Payer as Address.PayerID = Payer.ID
func (a *Address) QCreateAddress(db *gorm.DB, p *Payer) (int, error) {
	var err error
	a.PayerID = p.ID
	a.CreatedAt = time.Now()
	if err = db.Raw(`INSERT INTO address(payer_id, state, city, zip_code, street, number, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		a.PayerID, a.State, a.City, a.ZipCode, a.Street, a.Number, a.CreatedAt).Scan(&a.ID).Error; err != nil {
		return 500, err
	}
	p.AddressID = a.ID
	return p.QUpdatePayer(db)
}

// QGetPayers - Get all Payers
func (p *Payer) QGetPayers(db *gorm.DB, start int, count int) ([]Payer, int, error) {
	var payers []Payer
	if err := db.Model(&Payer{}).Preload("Address").Limit(count).Offset(start).Find(&payers).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return payers, 404, err
		default:
			return payers, 500, err
		}
	}
	return payers, 200, nil
}

// QGetPayer - Get Payer by ID
func (p *Payer) QGetPayer(db *gorm.DB) (int, error) {
	if err := db.Preload("Address").Where("payer.id=?", p.ID).First(&p).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}

func (p *Payer) QUpdatePayer(db *gorm.DB) (int, error) {
	var err error
	if err = validator.Validate(p); err != nil {
		return 400, err
	}

	if err = validator.Validate(p.Address); err != nil {
		return 400, err
	}

	p.UpdatedAt = time.Now()
	if err = db.Model(&p).Updates(p).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}

func (p *Payer) QPrimaryCard(db *gorm.DB, card_id int) (int, error) {
	var err error
	card := Card{ID: card_id}
	if code, err := card.QGetCard(db); err != nil {
		return code, err
	}
	if card.PayerID != p.ID {
		return 400, errors.New("invalid card id")
	}
	if err = db.Model(&p).Update("card_id", card_id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return 404, err
		default:
			return 500, err
		}
	}
	return 200, nil
}
