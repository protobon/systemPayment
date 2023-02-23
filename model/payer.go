package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Payer example
type Payer struct {
	ID            int            `json:"id" gorm:"primaryKey" example:"1" swaggerignore:"true"`
	Name          *string        `json:"name" example:"Jhon Doe" validate:"nonzero,min=6,max=100"`
	Email         *string        `json:"email" example:"jhondoe@mail.com" validate:"nonzero,min=6,max=100"`
	BirthDate     *string        `json:"birth_date" example:"24/07/1992" validate:"nonzero"`
	Phone         *string        `json:"phone" example:"+123456789" validate:"nonzero"`
	Document      *string        `json:"document" example:"23415162" validate:"nonzero"`
	UserReference *string        `json:"user_reference" example:"12345" validate:"nonzero"`
	Address       Address        `json:"address" gorm:"foreignKey:PayerID;references:ID" validate:"nonzero"`
	AddressID     int            `json:"-" swaggerignore:"true"`
	CardID        int            `json:"-" swaggerignore:"true"`
	Card          Card           `json:"card" gorm:"foreignKey:PayerID;references:ID" validate:"nonzero"`
	CreatedAt     time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt     time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt     gorm.DeletedAt `json:"-" swaggerignore:"true"`
}

type Address struct {
	ID        int            `json:"-" gorm:"primaryKey" example:"1" swaggerignore:"true"`
	PayerID   int            `json:"-" gorm:"column:payer_id" example:"1" swaggerignore:"true"`
	State     *string        `json:"state" example:"Rio de Janeiro" validate:"nonzero"`
	City      *string        `json:"city" example:"Volta Redonda" validate:"nonzero"`
	ZipCode   *string        `json:"zip_code" example:"27275-595" validate:"nonzero"`
	Street    *string        `json:"street" example:"Servidão B-1" validate:"nonzero"`
	Number    *string        `json:"number" example:"1106" validate:"nonzero"`
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"-" swaggerignore:"true"`
}

func (Payer) TableName() string {
	return "payer"
}

func (Address) TableName() string {
	return "address"
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
	if err = db.Raw(`INSERT INTO payer(name, email, birth_date, phone, document, user_reference, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`, p.Name, p.Email, p.BirthDate, p.Phone,
		p.Document, p.UserReference, p.CreatedAt).Scan(&p.ID).Error; err != nil {
		return 500, err
	}

	// Insert Payer's address
	code, err := p.Address.QCreateAddress(db, p)
	return code, err
}

// QCreateAddress - Insert into address
//
// Inserts new Address linked to Payer as Address.PayerID = Payer.ID
func (a *Address) QCreateAddress(db *gorm.DB, p *Payer) (int, error) {
	var err error
	a.PayerID = p.ID
	a.CreatedAt = time.Now()
	// err = db.Create(a).Error
	if err = db.Raw(`INSERT INTO address(payer_id, state, city, zip_code, street, number, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		a.PayerID, a.State, a.City, a.ZipCode, a.Street, a.Number, a.CreatedAt).Scan(&a.ID).Error; err != nil {
		return 500, err
	}
	p.AddressID = a.ID
	code, err := p.QUpdatePayer(db)
	return code, err
}

// QGetPayers - Get all Payers
func (p *Payer) QGetPayers(db *gorm.DB, start int, count int) ([]Payer, int, error) {
	var payers []Payer
	rows, err := db.Raw(`SELECT payer.id, payer.name, payer.email, payer.birth_date, payer.phone, payer.document, 
	payer.user_reference, payer.created_at, address.state, address.city, address.zip_code, address.street, 
	address.number, address.created_at FROM payer LEFT JOIN address ON address.id=payer.address_id LIMIT ? OFFSET ?`,
		count, start).Rows()
	if err != nil {
		return payers, 500, err
	}
	defer rows.Close()
	for rows.Next() {
		payer := Payer{}
		if err = rows.Scan(&payer.ID, &payer.Name, &payer.Email, &payer.BirthDate, &payer.Phone,
			&payer.Document, &payer.UserReference, &payer.CreatedAt, &payer.Address.State, &payer.Address.City,
			&payer.Address.ZipCode, &payer.Address.Street, &payer.Address.Number, &payer.Address.CreatedAt); err != nil {
			return payers, 500, err
		}
		payers = append(payers, payer)
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
