package model

import (
	"database/sql"
	"log"
	"time"
)

type DummyObject struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *DummyObject) QGetDummy(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM dummy WHERE id=$1",
		d.ID).Scan(&d.ID, &d.Name, &d.Price, &d.CreatedAt, &d.UpdatedAt)
}

func (d *DummyObject) QUpdateDummy(db *sql.DB) error {
	d.UpdatedAt = time.Now()
	_, err :=
		db.Exec("UPDATE dummy SET name=$1, price=$2, updated_at=$3 WHERE id=$4",
			d.Name, d.Price, d.UpdatedAt, d.ID)
	return err
}

func (d *DummyObject) QDeleteDummy(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dummy WHERE id=$1", d.ID)
	return err
}

func (d *DummyObject) QCreateDummy(db *sql.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	err := db.QueryRow(
		"INSERT INTO dummy(name, price, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id",
		d.Name, d.Price, d.CreatedAt, d.UpdatedAt).Scan(&d.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DummyObject) QGetDummies(db *sql.DB, start int, count int) ([]DummyObject, error) {
	rows, err := db.Query(
		"SELECT * FROM dummy LIMIT $1 OFFSET $2",
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

	var dummies []DummyObject

	for rows.Next() {
		var dummy DummyObject
		if err = rows.Scan(&dummy.ID, &dummy.Name, &dummy.Price,
			&dummy.CreatedAt, &dummy.UpdatedAt); err != nil {
			return nil, err
		}
		dummies = append(dummies, dummy)
	}

	return dummies, nil
}
