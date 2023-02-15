package model

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"gopkg.in/validator.v2"
)

type AddressObject struct {
	ID        int       `json:"id"`
	Payer     int       `json:"payer"`
	State     *string   `json:"state" validate:"nonzero"`
	City      *string   `json:"city" validate:"nonzero"`
	ZipCode   *string   `json:"zip_code" validate:"nonzero"`
	Street    *string   `json:"street" validate:"nonzero"`
	Number    *string   `json:"number" validate:"nonzero"`
	CreatedAt time.Time `json:"created_at"`
}

type DlocalAddress struct {
	State   *string `json:"state" validate:"nonzero"`
	City    *string `json:"city" validate:"nonzero"`
	ZipCode *string `json:"zip_code" validate:"nonzero"`
	Street  *string `json:"street" validate:"nonzero"`
	Number  *string `json:"number" validate:"nonzero"`
}

type PayerObject struct {
	ID            int           `json:"id"`
	Name          *string       `json:"name" validate:"nonzero"`
	Email         *string       `json:"email" validate:"nonzero"`
	BirthDate     *string       `json:"birth_date" validate:"nonzero"`
	Phone         *string       `json:"phone" validate:"nonzero"`
	Document      *string       `json:"document" validate:"nonzero"`
	UserReference *string       `json:"user_reference" validate:"nonzero"`
	Address       AddressObject `json:"address" validate:"nonzero"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type DlocalPayer struct {
	Name          string        `json:"name" validate:"nonzero"`
	Email         string        `json:"email" validate:"nonzero"`
	Document      string        `json:"document" validate:"nonzero"`
	UserReference string        `json:"user_reference" validate:"nonzero"`
	Address       AddressObject `json:"address" validate:"nonzero"`
	IP            string        `json:"ip" validate:"nonzero"`
	DeviceID      string        `json:"device_id" validate:"nonzero"`
}

func (p *PayerObject) QGetPayerFromID(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM payer WHERE id=$1",
		p.ID).Scan(&p.ID)
}

func (p *PayerObject) QGetPayerFromEmail(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM payer WHERE email=$1",
		p.Email).Scan(&p.ID)
}

func (p *PayerObject) QCreatePayer(db *sql.DB) error {
	if err := p.QGetPayerFromEmail(db); err == nil {
		// Found Payer from Email
		return errors.New("email already taken")
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	err := db.QueryRow(
		`INSERT INTO payer(name, email, birth_date, phone, document, user_reference, created_at, updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, p.Name, p.Email, p.BirthDate, p.Phone,
		p.Document, p.UserReference, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)

	if err != nil {
		return err
	}

	if err = validator.Validate(p.Address); err != nil {
		return err
	}

	p.Address.Payer = p.ID
	if err = p.Address.QCreateAddress(db); err != nil {
		return err
	}
	return nil
}

func (a *AddressObject) QCreateAddress(db *sql.DB) error {
	var err error
	a.CreatedAt = time.Now()
	err = db.QueryRow(
		`INSERT INTO address(payer, state, city, zip_code, street, number, created_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		a.Payer, a.State, a.City, a.ZipCode, a.Street, a.Number, a.CreatedAt).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *PayerObject) QGetPayer(db *sql.DB) error {
	return db.QueryRow(`SELECT payer.id, payer.name, payer.email, payer.birth_date, payer.phone, payer.document, 
		payer.user_reference, payer.created_at, address.state, address.city, address.zip_code, address.street, 
		address.number FROM payer JOIN address ON address.payer=payer.id WHERE payer.id=$1`, p.ID).Scan(
		&p.Name, &p.Email, &p.BirthDate, &p.Phone, &p.Document, &p.UserReference, &p.CreatedAt,
		&p.Address.State, &p.Address.City, &p.Address.ZipCode, &p.Address.Street, &p.Address.Number)
}

func (p *PayerObject) QGetPayers(db *sql.DB, start int, count int) ([]PayerObject, error) {
	rows, err := db.Query(
		`SELECT payer.id, payer.name, payer.email, payer.birth_date, payer.phone, payer.document, 
		payer.user_reference, payer.created_at, address.state, address.city, address.zip_code, address.street, 
		address.number FROM payer LEFT JOIN address ON address.payer=payer.id LIMIT $1 OFFSET $2`,
		count, start)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	var payers []PayerObject

	for rows.Next() {
		var payer PayerObject
		if err = rows.Scan(&payer.ID, &payer.Name, &payer.Email, &payer.BirthDate, &payer.Phone,
			&payer.Document, &payer.UserReference, &payer.CreatedAt, &payer.Address.State,
			&payer.Address.City, &payer.Address.ZipCode, &payer.Address.Street, &payer.Address.Number); err != nil {
			return nil, err
		}
		payers = append(payers, payer)
	}

	return payers, nil
}

func (p *PayerObject) QUpdatePayer(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM payer WHERE id=$1",
		p.ID).Scan(&p.ID)
}

func (p *PayerObject) QDeletePayer(db *sql.DB) error {
	return db.QueryRow("SELECT id FROM payer WHERE id=$1",
		p.ID).Scan(&p.ID)
}
