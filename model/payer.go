package model

import (
	"time"

	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Payer example
type Payer struct {
	ID            int            `gorm:"primaryKey" example:"1"`
	Name          *string        `example:"Jhon Doe" validate:"nonzero,min=6,max=100"`
	Email         *string        `example:"jhondoe@mail.com" validate:"nonzero,min=6,max=100"`
	BirthDate     *string        `example:"24/07/1992" validate:"nonzero"`
	Phone         *string        `example:"+123456789" validate:"nonzero"`
	Document      *string        `example:"23415162" validate:"nonzero"`
	UserReference *string        `example:"12345" validate:"nonzero"`
	Address       Address        `gorm:"foreignKey:PayerID;references:ID" validate:"nonzero"`
	AddressID     int            `json:"-" swaggerignore:"true"`
	CreatedAt     time.Time      `json:"-" swaggerignore:"true"`
	UpdatedAt     time.Time      `json:"-" swaggerignore:"true"`
	DeletedAt     gorm.DeletedAt `json:"-" swaggerignore:"true"`
}

type Address struct {
	ID        int            `json:"-" gorm:"primaryKey" example:"1" swaggerignore:"true"`
	PayerID   int            `json:"-" gorm:"column:payer_id" example:"1" swaggerignore:"true"`
	State     *string        `example:"Rio de Janeiro" validate:"nonzero"`
	City      *string        `example:"Volta Redonda" validate:"nonzero"`
	ZipCode   *string        `example:"27275-595" validate:"nonzero"`
	Street    *string        `example:"Servidão B-1" validate:"nonzero"`
	Number    *string        `example:"1106" validate:"nonzero"`
	CreatedAt time.Time      `json:"-" swaggerignore:"true"`
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
func (p *Payer) QCreatePayer(db *gorm.DB) (error, int) {
	var err error
	if err = validator.Validate(p); err != nil {
		return err, 400
	}
	if err = validator.Validate(p.Address); err != nil {
		return err, 400
	}

	p.CreatedAt = time.Now()
	// Create Payer
	if err = db.Raw(`INSERT INTO payer(name, email, birth_date, phone, document, user_reference, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`, p.Name, p.Email, p.BirthDate, p.Phone,
		p.Document, p.UserReference, p.CreatedAt).Scan(&p.ID).Error; err != nil {
		return err, 500
	}

	// Insert Payer's address
	err, code := p.Address.QCreateAddress(db, p)
	return err, code
}

// QCreateAddress - Insert into address
//
// Inserts new Address linked to Payer as Address.PayerID = Payer.ID
func (a *Address) QCreateAddress(db *gorm.DB, p *Payer) (error, int) {
	var err error
	a.PayerID = p.ID
	a.CreatedAt = time.Now()
	// err = db.Create(a).Error
	if err = db.Raw(`INSERT INTO address(payer_id, state, city, zip_code, street, number, created_at) 
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		a.PayerID, a.State, a.City, a.ZipCode, a.Street, a.Number, a.CreatedAt).Scan(&a.ID).Error; err != nil {
		return err, 500
	}
	p.AddressID = a.ID
	err, code := p.QUpdatePayer(db)
	return err, code
}

func (p *Payer) QGetPayers(db *gorm.DB, start int, count int) ([]Payer, error, int) {
	var payers []Payer
	rows, err := db.Raw(`SELECT payer.id, payer.name, payer.email, payer.birth_date, payer.phone, payer.document, 
	payer.user_reference, payer.created_at, address.state, address.city, address.zip_code, address.street, 
	address.number FROM payer LEFT JOIN address ON address.id=payer.address_id LIMIT ? OFFSET ?`,
		count, start).Rows()
	if err != nil {
		return payers, err, 500
	}
	defer rows.Close()
	for rows.Next() {
		payer := Payer{}
		if err = rows.Scan(&payer.ID, &payer.Name, &payer.Email, &payer.BirthDate, &payer.Phone,
			&payer.Document, &payer.UserReference, &payer.CreatedAt, &payer.Address.State,
			&payer.Address.City, &payer.Address.ZipCode, &payer.Address.Street, &payer.Address.Number); err != nil {
			return payers, err, 500
		}
		payers = append(payers, payer)
	}
	return payers, nil, 200
}

func (p *Payer) QGetPayer(db *gorm.DB) error {
	err := db.Preload("Address").Where("id = ? AND address.id = ?", p.ID, p.AddressID).First(&p).Error
	return err
}

func (p *Payer) QUpdatePayer(db *gorm.DB) (error, int) {
	var err error
	if err = validator.Validate(p); err != nil {
		return err, 400
	}

	if err = validator.Validate(p.Address); err != nil {
		return err, 400
	}

	p.UpdatedAt = time.Now()
	if err = db.Model(&p).Updates(p).Error; err != nil {
		return err, 500
	}
	return nil, 200
}
