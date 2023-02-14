package api

import (
	"database/sql"
	"fmt"
	"log"
	"systemPayment/api/routes"
	"systemPayment/database"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *App) Run(addr string) {
	router := a.Router

	err := router.Run(addr)
	if err != nil {
		return
	}
}

func (a *App) Initialize(user string, password string,
	dbhost string, dbname string) {
	fmt.Println("Initializing App...")
	var err error
	a.DB = database.DBInit(user, password, dbhost, dbname)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = gin.Default()

	var dummy = routes.Dummy{Router: a.Router, DB: a.DB}

	dummy.InitializeRoutes(dummy.DB)
	// schedule.RunCronJobs(a.DB)
	fmt.Println("***** App Running *****")
}
