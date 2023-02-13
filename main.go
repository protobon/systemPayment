package main

import (
	"os"
	"systemPayment/api"
)

func main() {
	app := api.App{}
	var user string
	var password string
	var dbname string
	var port string

	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname = os.Getenv("POSTGRES_DB")
	port = os.Getenv("APPLICATION_PORT")

	app.Initialize(
		user,
		password,
		dbname,
	)

	app.Run(port)
}
