package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func DBInit(user string, password string,
	dbhost string, dbname string) *sql.DB {
	connectionString :=
		fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
			user,
			password,
			dbhost,
			dbname)

	var err error
	var db *sql.DB
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	CreateTableDummy(db)
	return db
}

func CreateTableDummy(db *sql.DB) {
	if _, err := db.Exec(DummyTableCreate); err != nil {
		log.Fatal(err)
	}
}
